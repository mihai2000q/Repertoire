import Song from '../../types/models/Song.ts'
import {
  ActionIcon,
  alpha,
  Avatar,
  Center,
  Flex,
  Group,
  Menu,
  NumberFormatter,
  Space,
  Stack,
  Text,
  Tooltip
} from '@mantine/core'
import dayjs from 'dayjs'
import { useAppDispatch } from '../../state/store.ts'
import { openAlbumDrawer, openSongDrawer } from '../../state/slice/globalSlice.ts'
import { MouseEvent } from 'react'
import { IconCircleMinus, IconDisc, IconDots, IconEye, IconTrash } from '@tabler/icons-react'
import { useDisclosure, useHover } from '@mantine/hooks'
import WarningModal from '../@ui/modal/WarningModal.tsx'
import Order from '../../types/Order.ts'
import SongProperty from '../../types/enums/SongProperty.ts'
import DifficultyBar from '../@ui/bar/DifficultyBar.tsx'
import ConfidenceBar from '../@ui/bar/ConfidenceBar.tsx'
import ProgressBar from '../@ui/bar/ProgressBar.tsx'
import { useNavigate } from 'react-router-dom'
import PerfectRehearsalMenuItem from '../@ui/menu/item/song/PerfectRehearsalMenuItem.tsx'
import PartialRehearsalMenuItem from '../@ui/menu/item/song/PartialRehearsalMenuItem.tsx'
import { useRemoveSongsFromArtistMutation } from '../../state/api/artistsApi.ts'
import { useDeleteSongMutation } from '../../state/api/songsApi.ts'
import CustomIconMusicNoteEighth from '../@ui/icons/CustomIconMusicNoteEighth.tsx'
import OpenLinksMenuItem from '../@ui/menu/item/song/OpenLinksMenuItem.tsx'
import YoutubeModal from '../@ui/modal/YoutubeModal.tsx'
import AddToPlaylistMenuItem from '../@ui/menu/item/AddToPlaylistMenuItem.tsx'
import { ContextMenu } from '../@ui/menu/ContextMenu.tsx'
import useDoubleMenu from '../../hooks/useDoubleMenu.ts'
import { toast } from 'react-toastify'

interface ArtistSongCardProps {
  song: Song
  artistId: string
  isUnknownArtist: boolean
  order: Order
}

function ArtistSongCard({ song, artistId, isUnknownArtist, order }: ArtistSongCardProps) {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const { ref, hovered } = useHover()

  const [removeSongsFromArtist, { isLoading: isRemoveLoading }] = useRemoveSongsFromArtistMutation()
  const [deleteSong, { isLoading: isDeleteLoading }] = useDeleteSongMutation()

  const { openedMenu, toggleMenu, openedContextMenu, toggleContextMenu, closeMenus } =
    useDoubleMenu()

  const isSelected = hovered || openedMenu || openedContextMenu

  const [openedYoutube, { open: openYoutube, close: closeYoutube }] = useDisclosure(false)
  const [openedRemoveWarning, { open: openRemoveWarning, close: closeRemoveWarning }] =
    useDisclosure(false)
  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  function handleClick() {
    dispatch(openSongDrawer(song.id))
  }

  function handleAlbumClick(e: MouseEvent) {
    e.stopPropagation()
    dispatch(openAlbumDrawer(song.album.id))
  }

  function handleViewDetails(e: MouseEvent) {
    e.stopPropagation()
    navigate(`/song/${song.id}`)
  }

  function handleViewAlbum() {
    navigate(`/album/${song.album.id}`)
  }

  function handleOpenRemoveWarning(e: MouseEvent) {
    e.stopPropagation()
    openRemoveWarning()
  }

  function handleOpenDeleteWarning(e: MouseEvent) {
    e.stopPropagation()
    openDeleteWarning()
  }

  async function handleRemoveFromArtist() {
    await removeSongsFromArtist({ songIds: [song.id], id: artistId }).unwrap()
    toast.success(`${song.title} removed from artist!`)
  }

  async function handleDelete() {
    await deleteSong(song.id).unwrap()
    toast.success(`${song.title} deleted!`)
  }

  const menuDropdown = (
    <>
      <Menu.Item leftSection={<IconEye size={14} />} onClick={handleViewDetails}>
        View Details
      </Menu.Item>
      <Menu.Item
        leftSection={<IconDisc size={14} />}
        disabled={!song.album}
        onClick={handleViewAlbum}
      >
        View Album
      </Menu.Item>
      <OpenLinksMenuItem song={song} openYoutube={openYoutube} />

      <Menu.Divider />
      <AddToPlaylistMenuItem ids={[song.id]} type={'song'} closeMenu={closeMenus} />
      <PartialRehearsalMenuItem songId={song.id} closeMenu={closeMenus} />
      <PerfectRehearsalMenuItem songId={song.id} closeMenu={closeMenus} />
      <Menu.Divider />

      {!isUnknownArtist && (
        <Menu.Item leftSection={<IconCircleMinus size={14} />} onClick={handleOpenRemoveWarning}>
          Remove from Artist
        </Menu.Item>
      )}
      <Menu.Item
        leftSection={<IconTrash size={14} />}
        c={'red.5'}
        onClick={handleOpenDeleteWarning}
      >
        Delete
      </Menu.Item>
    </>
  )

  return (
    <ContextMenu shadow={'lg'} opened={openedContextMenu} onChange={toggleContextMenu}>
      <ContextMenu.Target>
        <Group
          ref={ref}
          aria-label={`song-card-${song.title}`}
          wrap={'nowrap'}
          sx={(theme) => ({
            cursor: 'default',
            transition: '0.3s',
            ...(isSelected && {
              boxShadow: theme.shadows.xl,
              backgroundColor: alpha(theme.colors.primary[0], 0.15)
            })
          })}
          px={'md'}
          py={'xs'}
          gap={0}
          onClick={handleClick}
        >
          <Avatar
            radius={'md'}
            src={song.imageUrl ?? song.album?.imageUrl}
            alt={(song.imageUrl ?? song.album?.imageUrl) && song.title}
            bg={'gray.5'}
          >
            <Center c={'white'}>
              <CustomIconMusicNoteEighth aria-label={`default-icon-${song.title}`} size={20} />
            </Center>
          </Avatar>

          <Space
            ml={{ base: 'xs', xs: 'md', sm: 'xs', betweenSmMd: 'md', md: 'sm', lg: 'md' }}
            style={{ transition: '0.16s' }}
          />

          <Stack gap={0} flex={1} style={{ overflow: 'hidden' }}>
            <Group gap={'0px 4px'}>
              <Text fw={500} truncate={'end'}>
                {song.title}
              </Text>
              {song.album && (
                <>
                  <Text fz={'sm'}>-</Text>
                  <Text
                    fz={'sm'}
                    c={'dimmed'}
                    truncate={'end'}
                    sx={{ '&:hover': { textDecoration: 'underline' } }}
                    style={{ cursor: 'pointer' }}
                    onClick={handleAlbumClick}
                    lh={'xxs'}
                  >
                    {song.album.title}
                  </Text>
                </>
              )}
            </Group>
            <Flex>
              {order.property === SongProperty.ReleaseDate && (
                <Tooltip
                  label={`Song was released on ${dayjs(song.releaseDate).format('D MMMM YYYY')}`}
                  openDelay={400}
                  disabled={!song.releaseDate}
                >
                  <Text fz={'xs'} c={'dimmed'}>
                    {song.releaseDate && dayjs(song.releaseDate).format('DD MMM YYYY')}
                  </Text>
                </Tooltip>
              )}
              {order.property === SongProperty.Difficulty && (
                <DifficultyBar difficulty={song.difficulty} mt={4} miw={110} />
              )}
              {order.property === SongProperty.Rehearsals && (
                <Tooltip.Floating
                  role={'tooltip'}
                  label={
                    <>
                      Rehearsals: <NumberFormatter value={song.rehearsals} />
                    </>
                  }
                >
                  <Text fz={'xs'} c={'dimmed'}>
                    <NumberFormatter value={song.rehearsals} />
                  </Text>
                </Tooltip.Floating>
              )}
              {order.property === SongProperty.Confidence && (
                <ConfidenceBar confidence={song.confidence} w={100} mt={4} />
              )}
              {order.property === SongProperty.Progress && (
                <ProgressBar progress={song.progress} w={100} mt={4} />
              )}
              {order.property === SongProperty.LastPlayed && (
                <Tooltip
                  label={`Song was played last time on ${dayjs(song.lastTimePlayed).format('D MMMM YYYY [at] hh:mm A')}`}
                  openDelay={400}
                  disabled={!song.lastTimePlayed}
                >
                  <Text fz={'xs'} c={'dimmed'}>
                    {song.lastTimePlayed
                      ? dayjs(song.lastTimePlayed).format('DD MMM YYYY')
                      : 'never'}
                  </Text>
                </Tooltip>
              )}
            </Flex>
          </Stack>

          <Space
            ml={{ base: 'xs', xs: 'md', sm: 'xs', betweenSmMd: 'md', md: 'sm', lg: 'md' }}
            style={{ transition: '0.16s' }}
          />

          <Menu opened={openedMenu} onChange={toggleMenu}>
            <Menu.Target>
              <ActionIcon
                size={'md'}
                variant={'grey'}
                aria-label={'more-menu'}
                onClick={(e) => e.stopPropagation()}
                style={{
                  transition: '0.3s',
                  opacity: isSelected ? 1 : 0
                }}
              >
                <IconDots size={15} />
              </ActionIcon>
            </Menu.Target>
            <Menu.Dropdown>{menuDropdown}</Menu.Dropdown>
          </Menu>
        </Group>
      </ContextMenu.Target>

      <ContextMenu.Dropdown>{menuDropdown}</ContextMenu.Dropdown>

      <YoutubeModal
        title={song.title}
        link={song.youtubeLink}
        opened={openedYoutube}
        onClose={closeYoutube}
      />
      <WarningModal
        opened={openedRemoveWarning}
        onClose={closeRemoveWarning}
        title={`Remove Song From Artist`}
        description={
          <Group gap={'xxs'}>
            <Text>Are you sure you want to remove</Text>
            <Text fw={600}>{song.title}</Text>
            <Text>from this artist?</Text>
          </Group>
        }
        isLoading={isRemoveLoading}
        onYes={handleRemoveFromArtist}
      />
      <WarningModal
        opened={openedDeleteWarning}
        onClose={closeDeleteWarning}
        title={`Delete Song`}
        description={
          <Group gap={'xxs'}>
            <Text>Are you sure you want to delete</Text>
            <Text fw={600}>{song.title}</Text>
            <Text>?</Text>
          </Group>
        }
        isLoading={isDeleteLoading}
        onYes={handleDelete}
      />
    </ContextMenu>
  )
}

export default ArtistSongCard
