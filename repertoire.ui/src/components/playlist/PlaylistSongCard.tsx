import Song from '../../types/models/Song.ts'
import {
  ActionIcon,
  alpha,
  Center,
  Flex,
  Grid,
  Group,
  Menu,
  NumberFormatter,
  Stack,
  Text,
  Tooltip
} from '@mantine/core'
import { useAppDispatch } from '../../state/store.ts'
import { openAlbumDrawer, openArtistDrawer, openSongDrawer } from '../../state/slice/globalSlice.ts'
import { useDisclosure, useHover, useMergedRef } from '@mantine/hooks'
import { forwardRef, MouseEvent } from 'react'
import { IconCircleMinus, IconDisc, IconDots, IconEye, IconUser } from '@tabler/icons-react'
import WarningModal from '../@ui/modal/WarningModal.tsx'
import { useNavigate } from 'react-router-dom'
import { DraggableProvided } from '@hello-pangea/dnd'
import PerfectRehearsalMenuItem from '../@ui/menu/item/PerfectRehearsalMenuItem.tsx'
import PartialRehearsalMenuItem from '../@ui/menu/item/song/PartialRehearsalMenuItem.tsx'
import CustomIconMusicNoteEighth from '../@ui/icons/CustomIconMusicNoteEighth.tsx'
import { useRemoveSongsFromPlaylistMutation } from '../../state/api/playlistsApi.ts'
import SongProperty from '../../types/enums/properties/SongProperty.ts'
import Order from '../../types/Order.ts'
import DifficultyBar from '../@ui/bar/DifficultyBar.tsx'
import ConfidenceBar from '../@ui/bar/ConfidenceBar.tsx'
import ProgressBar from '../@ui/bar/ProgressBar.tsx'
import dayjs from 'dayjs'
import YoutubeModal from '../@ui/modal/YoutubeModal.tsx'
import OpenLinksMenuItem from '../@ui/menu/item/song/OpenLinksMenuItem.tsx'
import AddToPlaylistMenuItem from '../@ui/menu/item/AddToPlaylistMenuItem.tsx'
import useDoubleMenu from '../../hooks/useDoubleMenu.ts'
import { ContextMenu } from '../@ui/menu/ContextMenu.tsx'
import { toast } from 'react-toastify'
import useClickSelectSelectable from '../../hooks/useClickSelectSelectable.ts'
import SelectableAvatar from '../@ui/image/SelectableAvatar.tsx'

interface PlaylistSongCardProps {
  song: Song
  playlistId: string
  order: Order
  isDragging: boolean
  draggableProvided?: DraggableProvided
}

const PlaylistSongCard = forwardRef<HTMLElement, PlaylistSongCardProps>(
  ({ song, playlistId, order, isDragging, draggableProvided }, refProp) => {
    const dispatch = useAppDispatch()
    const navigate = useNavigate()
    const {
      ref: selectableRef,
      isClickSelected,
      isClickSelectionActive,
      isLastInSelection
    } = useClickSelectSelectable(song.playlistSongId)
    const { ref: hoverRef, hovered } = useHover()
    const ref = useMergedRef(hoverRef, draggableProvided?.innerRef, refProp, selectableRef)

    const [removeSongsFromPlaylist, { isLoading: isRemoveLoading }] =
      useRemoveSongsFromPlaylistMutation()

    const [openedYoutube, { open: openYoutube, close: closeYoutube }] = useDisclosure(false)

    const { openedMenu, toggleMenu, openedContextMenu, toggleContextMenu, closeMenus } =
      useDoubleMenu()

    const isSelected = hovered || openedContextMenu || openedMenu || isDragging || isClickSelected

    const [openedRemoveWarning, { open: openRemoveWarning, close: closeRemoveWarning }] =
      useDisclosure(false)

    function handleClick(e: MouseEvent) {
      if (e.ctrlKey || e.shiftKey) return
      dispatch(openSongDrawer(song.id))
    }

    function handleAlbumClick(e: MouseEvent) {
      if (e.ctrlKey || e.shiftKey) return
      e.stopPropagation()
      dispatch(openAlbumDrawer(song.album.id))
    }

    function handleArtistClick(e: MouseEvent) {
      if (e.ctrlKey || e.shiftKey) return
      e.stopPropagation()
      dispatch(openArtistDrawer(song.artist.id))
    }

    function handleViewDetails(e: MouseEvent) {
      e.stopPropagation()
      navigate(`/song/${song.id}`)
    }

    function handleViewArtist(e: MouseEvent) {
      e.stopPropagation()
      navigate(`/artist/${song.artist.id}`)
    }

    function handleViewAlbum(e: MouseEvent) {
      e.stopPropagation()
      navigate(`/album/${song.album.id}`)
    }

    function handleOpenRemoveWarning(e: MouseEvent) {
      e.stopPropagation()
      openRemoveWarning()
    }

    async function handleRemoveFromPlaylist() {
      await removeSongsFromPlaylist({
        playlistSongIds: [song.playlistSongId],
        id: playlistId
      }).unwrap()
      toast.success(`${song.title} removed from playlist!`)
    }

    const menuDropdown = (
      <>
        <Menu.Item leftSection={<IconEye size={14} />} onClick={handleViewDetails}>
          View Details
        </Menu.Item>
        <Menu.Item
          leftSection={<IconUser size={14} />}
          disabled={!song.artist}
          onClick={handleViewArtist}
        >
          View Artist
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
        <AddToPlaylistMenuItem ids={[song.id]} type={'songs'} closeMenu={closeMenus} />
        <PartialRehearsalMenuItem songId={song.id} closeMenu={closeMenus} />
        <PerfectRehearsalMenuItem id={song.id} closeMenu={closeMenus} type={'song'} />
        <Menu.Divider />

        <Menu.Item
          leftSection={<IconCircleMinus size={14} />}
          c={'red.5'}
          onClick={handleOpenRemoveWarning}
        >
          Remove from Playlist
        </Menu.Item>
      </>
    )

    return (
      <ContextMenu
        opened={openedContextMenu}
        onChange={toggleContextMenu}
        disabled={isClickSelectionActive}
      >
        <ContextMenu.Target>
          <Group
            ref={ref}
            wrap={'nowrap'}
            aria-label={`song-card-${song.title}`}
            aria-selected={isSelected}
            {...draggableProvided?.draggableProps}
            {...draggableProvided?.dragHandleProps}
            style={{
              ...draggableProvided?.draggableProps?.style,
              cursor: 'default'
            }}
            sx={(theme) => ({
              transition: '0.3s',
              borderRadius: 0,
              border: '1px solid transparent',
              ...(isSelected && {
                boxShadow: theme.shadows.md,
                backgroundColor: alpha(theme.colors.primary[0], 0.15)
              }),

              ...(isClickSelected && {
                boxShadow: 'none',
                backgroundColor: alpha(theme.colors.primary[0], 0.15),
                ...(hovered && {
                  boxShadow: theme.shadows.xs,
                  backgroundColor: alpha(theme.colors.primary[0], 0.35)
                }),
                ...(isLastInSelection && {
                  boxShadow: theme.shadows.lg
                })
              }),

              ...(isDragging && {
                boxShadow: theme.shadows.xl,
                borderRadius: '16px',
                backgroundColor: alpha(theme.white, 0.33),
                border: `1px solid ${alpha(theme.colors.primary[9], 0.33)}`
              })
            })}
            px={'md'}
            py={'xs'}
            onClick={handleClick}
          >
            <Grid columns={12} align={'center'} w={'100%'}>
              <Grid.Col
                span={
                  order.property === SongProperty.PlaylistTrackNo ||
                  order.property === SongProperty.Title
                    ? 'auto'
                    : 6
                }
              >
                <Group>
                  <Text fw={500} miw={30} maw={30} ta={'center'}>
                    {song.playlistTrackNo}
                  </Text>

                  <SelectableAvatar
                    isSelected={isClickSelected}
                    radius={'md'}
                    src={song.imageUrl ?? song.album?.imageUrl}
                    alt={(song.imageUrl ?? song.album?.imageUrl) && song.title}
                    bg={'gray.5'}
                  >
                    <Center c={'white'}>
                      <CustomIconMusicNoteEighth
                        aria-label={`default-icon-${song.title}`}
                        size={20}
                      />
                    </Center>
                  </SelectableAvatar>

                  <Stack flex={1} gap={0} style={{ overflow: 'hidden' }}>
                    <Group gap={'xxs'} wrap={'nowrap'}>
                      <Text fw={500} lineClamp={1}>
                        {song.title}
                      </Text>
                      {song.album && (
                        <>
                          <Text fz={'sm'}>-</Text>
                          <Text
                            fz={'sm'}
                            c={'dimmed'}
                            lineClamp={1}
                            sx={{ '&:hover': { textDecoration: 'underline' } }}
                            style={{ cursor: 'pointer' }}
                            onClick={handleAlbumClick}
                          >
                            {song.album.title}
                          </Text>
                        </>
                      )}
                    </Group>
                    {song.artist && (
                      <Text
                        fz={'sm'}
                        c={'dimmed'}
                        sx={{ '&:hover': { textDecoration: 'underline' } }}
                        style={{ cursor: 'pointer', alignSelf: 'start' }}
                        onClick={handleArtistClick}
                        lineClamp={1}
                      >
                        {song.artist.name}
                      </Text>
                    )}
                  </Stack>
                </Group>
              </Grid.Col>

              <Grid.Col
                span={
                  order.property === SongProperty.PlaylistTrackNo ||
                  order.property === SongProperty.Title
                    ? 0
                    : 'auto'
                }
              >
                <Flex px={'10%'}>
                  {order.property === SongProperty.ReleaseDate && (
                    <Tooltip
                      label={`Song was released on ${dayjs(song.releaseDate).format('D MMMM YYYY')}`}
                      openDelay={400}
                      disabled={!song.releaseDate}
                    >
                      <Text fw={500} c={'dimmed'} inline>
                        {song.releaseDate
                          ? dayjs(song.releaseDate).format('DD MMM YYYY')
                          : 'unknown'}
                      </Text>
                    </Tooltip>
                  )}
                  {order.property === SongProperty.Difficulty && (
                    <DifficultyBar difficulty={song.difficulty} miw={'max(15vw, 120px)'} />
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
                      <Text fw={500} c={'dimmed'} inline>
                        <NumberFormatter value={song.rehearsals} />
                      </Text>
                    </Tooltip.Floating>
                  )}
                  {order.property === SongProperty.Confidence && (
                    <ConfidenceBar confidence={song.confidence} flex={1} />
                  )}
                  {order.property === SongProperty.Progress && (
                    <ProgressBar progress={song.progress} flex={1} />
                  )}
                  {order.property === SongProperty.LastPlayed && (
                    <Tooltip
                      label={`Song was played last time on ${dayjs(song.lastTimePlayed).format('D MMMM YYYY [at] hh:mm A')}`}
                      openDelay={400}
                      disabled={!song.lastTimePlayed}
                    >
                      <Text fw={500} c={'dimmed'} inline>
                        {song.lastTimePlayed
                          ? dayjs(song.lastTimePlayed).format('DD MMM YYYY')
                          : 'never'}
                      </Text>
                    </Tooltip>
                  )}
                </Flex>
              </Grid.Col>

              <Grid.Col span={'content'}>
                <Menu opened={openedMenu} onChange={toggleMenu}>
                  <Menu.Target>
                    <ActionIcon
                      size={'md'}
                      variant={'grey'}
                      aria-label={'more-menu'}
                      disabled={isClickSelectionActive}
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
              </Grid.Col>
            </Grid>
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
          title={`Remove Song From Playlist`}
          description={
            <Group gap={'xxs'}>
              <Text>Are you sure you want to remove</Text>
              <Text fw={600}>{song.title}</Text>
              <Text>from this playlist?</Text>
            </Group>
          }
          isLoading={isRemoveLoading}
          onYes={handleRemoveFromPlaylist}
        />
      </ContextMenu>
    )
  }
)

PlaylistSongCard.displayName = 'PlaylistSongCard'

export default PlaylistSongCard
