import {
  alpha,
  Avatar,
  Button,
  Center,
  Divider,
  Group,
  Highlight,
  LoadingOverlay,
  Menu,
  Modal,
  ScrollArea,
  Stack,
  Text,
  TextInput
} from '@mantine/core'
import { IconPlaylist, IconPlaylistAdd, IconSearch } from '@tabler/icons-react'
import {
  useAddAlbumsToPlaylistMutation,
  useAddArtistsToPlaylistMutation,
  useAddSongsToPlaylistMutation,
  useGetPlaylistsQuery
} from '../../../../state/api/playlistsApi.ts'
import Playlist from '../../../../types/models/Playlist.ts'
import { useDebouncedValue, useDidUpdate, useSessionStorage } from '@mantine/hooks'
import useOrderBy from '../../../../hooks/api/useOrderBy.ts'
import OrderType from '../../../../types/enums/OrderType.ts'
import useFilters from '../../../../hooks/filter/useFilters.ts'
import FilterOperator from '../../../../types/enums/FilterOperator.ts'
import useFiltersHandlers from '../../../../hooks/filter/useFiltersHandlers.ts'
import useSearchBy from '../../../../hooks/api/useSearchBy.ts'
import PlaylistProperty from '../../../../types/enums/PlaylistProperty.ts'
import SessionStorageKeys from '../../../../types/enums/SessionStorageKeys.ts'
import plural from '../../../../utils/plural.ts'
import {
  AddAlbumsToPlaylistResponse,
  AddArtistsToPlaylistResponse,
  AddSongsToPlaylistResponse,
  AddToPlaylistResponse
} from '../../../../types/responses/PlaylistResponses.ts'
import { toast } from 'react-toastify'
import { ReactNode, useState } from 'react'

const AlreadyAddedModal = ({
  opened,
  onClose,
  text,
  withCancel,
  onRetry
}: {
  opened: boolean
  onClose: () => void
  text: ReactNode
  withCancel: boolean
  onRetry: (forceAdd: boolean) => void
}) => (
  <Modal
    opened={opened}
    onClose={onClose}
    title={'Already Added'}
    size={'auto'}
    zIndex={10000}
    centered
  >
    <Stack px={'xs'} align={'center'} maw={'max(min(40vw, 720px), 350px)'}>
      <Text lineClamp={1}>{text}</Text>
      <Group wrap={'nowrap'} gap={'xxs'}>
        <Button variant={'transparent'} onClick={withCancel ? onClose : () => onRetry(true)}>
          {withCancel ? 'Cancel' : 'Add All'}
        </Button>
        <Button variant={'filled'} onClick={() => onRetry(withCancel)}>
          {withCancel ? 'Add Anyway' : 'Just New Ones'}
        </Button>
      </Group>
    </Stack>
  </Modal>
)

function PlaylistOption({
  playlist,
  searchValue,
  onClick
}: {
  playlist: Playlist
  searchValue: string
  onClick: () => void
}) {
  return (
    <Group
      wrap={'nowrap'}
      role={'menuitem'}
      aria-label={playlist.title}
      py={6}
      px={8}
      gap={'xs'}
      sx={(theme) => ({
        cursor: 'pointer',
        borderRadius: theme.radius.md,
        backgroundColor: 'transparent',
        transition: '0.16s',
        '&:hover': {
          backgroundColor: alpha(theme.colors.gray[1], 0.7)
        }
      })}
      onClick={onClick}
    >
      <Avatar
        size={'sm'}
        radius={'md'}
        src={playlist.imageUrl}
        alt={playlist.imageUrl && playlist.title}
        bg={'gray.5'}
        style={(theme) => ({ aspectRatio: 1, boxShadow: theme.shadows.sm })}
      >
        <Center c={'white'}>
          <IconPlaylist aria-label={`default-icon-${playlist.title}`} size={12} />
        </Center>
      </Avatar>

      <Stack flex={1} gap={0}>
        <Highlight highlight={searchValue} fw={500} fz={'sm'} lh={'xs'} lineClamp={1}>
          {playlist.title}
        </Highlight>
        <Text fz={'xs'} c={'dimmed'} lh={'xxs'}>
          {playlist.songsCount} song{plural(playlist.songsCount)}
        </Text>
      </Stack>
    </Group>
  )
}

interface AddToPlaylistMenuItemProps {
  ids: string[]
  type: 'song' | 'album' | 'artist'
  closeMenu: () => void
  disabled?: boolean
}

function AddToPlaylistMenuItem({ ids, type, closeMenu, disabled }: AddToPlaylistMenuItemProps) {
  const [search, setSearch] = useSessionStorage({
    key: SessionStorageKeys.AddToPlaylist,
    defaultValue: ''
  })
  const [searchValue] = useDebouncedValue(search, 200)

  const orderBy = useOrderBy([
    { property: PlaylistProperty.LastModified, type: OrderType.Descending }
  ])
  const [filters, setFilters] = useFilters([
    { property: PlaylistProperty.Title, operator: FilterOperator.PatternMatching, isSet: false }
  ])
  const activeFilters = Array.from(filters.values()).filter((filter) => filter.isSet).length
  const { handleValueChange } = useFiltersHandlers(filters, setFilters)
  const searchBy = useSearchBy(filters)

  useDidUpdate(
    () =>
      handleValueChange(
        PlaylistProperty.Title + FilterOperator.PatternMatching,
        searchValue.trim()
      ),
    [searchValue]
  )

  const {
    data: playlists,
    isLoading,
    isFetching
  } = useGetPlaylistsQuery({
    orderBy: orderBy,
    searchBy: searchBy
  })

  const [addArtistsToPlaylist] = useAddArtistsToPlaylistMutation()
  const [addAlbumsToPlaylist] = useAddAlbumsToPlaylistMutation()
  const [addSongsToPlaylist] = useAddSongsToPlaylistMutation()

  const [alreadyAddedModalState, setAlreadyAddedModalState] = useState<{
    opened: boolean
    text: ReactNode
    withCancel: boolean
    onRetry: (forceAdd: boolean) => void
  }>({
    opened: false,
    text: '',
    withCancel: false,
    onRetry: () => {}
  })

  function closeWarning() {
    setAlreadyAddedModalState({ ...alreadyAddedModalState, opened: false })
    closeMenu()
  }

  async function handleClick(playlist: Playlist) {
    // send first request
    let res: AddToPlaylistResponse
    let added: string[]

    switch (type) {
      case 'artist':
        res = await addArtistsToPlaylist({ id: playlist.id, artistIds: ids }).unwrap()
        added = (res as AddArtistsToPlaylistResponse).addedSongIds
        break
      case 'album':
        res = await addAlbumsToPlaylist({ id: playlist.id, albumIds: ids }).unwrap()
        added = (res as AddAlbumsToPlaylistResponse).addedSongIds
        break
      case 'song':
        res = await addSongsToPlaylist({ id: playlist.id, songIds: ids }).unwrap()
        added = (res as AddSongsToPlaylistResponse).added
        break
    }

    // check success
    if (res.success) {
      toast.success(`Successfully added ${added.length} song${plural(added)}!`)
      closeMenu()
      return
    }

    // if it failed open already added modal
    openAlreadyAddedModal(res, playlist)
  }

  function openAlreadyAddedModal(res: AddToPlaylistResponse, playlist: Playlist) {
    let partialText: string
    let withCancel = true
    let retryFn: (forceAdd: boolean) => void

    switch (type) {
      case 'artist': {
        const { duplicateArtistIds } = res as AddArtistsToPlaylistResponse
        if (duplicateArtistIds.length > 1 && ids.length === duplicateArtistIds.length) {
          partialText = 'These artists are'
        } else if (duplicateArtistIds.length > 1) {
          partialText = 'Some artists are'
          withCancel = false
        } else if (duplicateArtistIds.length === 1) {
          partialText = 'This artist is'
        } else {
          partialText = 'Some songs are'
          withCancel = false
        }

        retryFn = async (forceAdd: boolean) => {
          const { addedSongIds } = await addArtistsToPlaylist({
            id: playlist.id,
            artistIds: ids,
            forceAdd
          }).unwrap()
          toast.success(`Successfully added ${addedSongIds.length} song${plural(addedSongIds)}!`)
        }
        break
      }

      case 'album': {
        const { duplicateAlbumIds } = res as AddAlbumsToPlaylistResponse
        if (duplicateAlbumIds.length > 1 && ids.length === duplicateAlbumIds.length) {
          partialText = 'These albums are'
        } else if (duplicateAlbumIds.length > 1) {
          partialText = 'Some albums are'
          withCancel = false
        } else if (duplicateAlbumIds.length === 1) {
          partialText = 'This album is'
        } else {
          partialText = 'Some songs are'
          withCancel = false
        }

        retryFn = async (forceAdd: boolean) => {
          const { addedSongIds } = await addAlbumsToPlaylist({
            id: playlist.id,
            albumIds: ids,
            forceAdd
          }).unwrap()
          toast.success(`Successfully added ${addedSongIds.length} song${plural(addedSongIds)}!`)
        }
        break
      }

      case 'song': {
        const { duplicates } = res as AddSongsToPlaylistResponse
        if (duplicates.length > 1 && ids.length === duplicates.length) {
          partialText = 'These songs are'
        } else if (duplicates.length > 1) {
          partialText = 'Some songs are'
          withCancel = false
        } else {
          partialText = 'This song is'
        }

        retryFn = async (forceAdd: boolean) => {
          const newRes = await addSongsToPlaylist({
            id: playlist.id,
            songIds: ids,
            forceAdd
          }).unwrap()
          toast.success(`Successfully added ${newRes.added.length} song${plural(newRes.added)}!`)
        }
        break
      }
    }

    const text = (
      <>
        {partialText} already in <b>{playlist.title}</b>
      </>
    )
    const onRetry = (forceAdd: boolean) => {
      retryFn(forceAdd)
      closeMenu()
    }

    setAlreadyAddedModalState({ opened: true, text, withCancel, onRetry })
  }

  return (
    <>
      <Menu.Sub>
        <Menu.Sub.Target>
          <Menu.Sub.Item
            leftSection={<IconPlaylistAdd size={14} />}
            disabled={
              disabled === true || isLoading || (playlists.totalCount === 0 && activeFilters === 0)
            }
            onClick={(e) => e.stopPropagation()} // TODO: Remove when mantine updates
          >
            Add To Playlist
          </Menu.Sub.Item>
        </Menu.Sub.Target>

        <Menu.Sub.Dropdown miw={150} maw={200} p={0}>
          <TextInput
            aria-label={'search'}
            variant={'unstyled'}
            size={'xs'}
            maxLength={100}
            placeholder={'Search'}
            leftSection={<IconSearch size={12} />}
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
          <Divider />

          <ScrollArea.Autosize mah={'max(250px, 50vh)'} scrollbars={'y'} scrollbarSize={5}>
            <Stack gap={'xxs'} py={7} px={'xxs'} style={{ transition: '0.16s' }}>
              <LoadingOverlay visible={isFetching} />
              {playlists?.models?.map((playlist) => (
                <PlaylistOption
                  key={playlist.id}
                  playlist={playlist}
                  searchValue={searchValue}
                  onClick={() => handleClick(playlist)}
                />
              ))}
              {playlists?.totalCount === 0 && activeFilters === 0 && (
                <Text fz={'xs'} c={'dimmed'} px={'xs'}>
                  There are no playlists
                </Text>
              )}
              {playlists?.totalCount === 0 && activeFilters > 0 && (
                <Text fz={'xs'} c={'dimmed'} px={'xs'}>
                  No playlists found
                </Text>
              )}
            </Stack>
          </ScrollArea.Autosize>
        </Menu.Sub.Dropdown>
      </Menu.Sub>

      <AlreadyAddedModal onClose={closeWarning} {...alreadyAddedModalState} />
    </>
  )
}

export default AddToPlaylistMenuItem
