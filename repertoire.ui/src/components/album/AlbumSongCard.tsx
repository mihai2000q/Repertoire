import Song from '../../types/models/Song.ts'
import {
  ActionIcon,
  alpha,
  Avatar,
  Center,
  Flex,
  Grid,
  Group,
  Menu,
  NumberFormatter,
  Text,
  Tooltip
} from '@mantine/core'
import { useAppDispatch } from '../../state/store.ts'
import { openSongDrawer } from '../../state/slice/globalSlice.ts'
import { useDisclosure, useHover, useMergedRef } from '@mantine/hooks'
import { MouseEvent } from 'react'
import { IconCircleMinus, IconDots, IconEye, IconTrash } from '@tabler/icons-react'
import WarningModal from '../@ui/modal/WarningModal.tsx'
import Order from '../../types/Order.ts'
import SongProperty from '../../types/enums/SongProperty.ts'
import ProgressBar from '../@ui/bar/ProgressBar.tsx'
import ConfidenceBar from '../@ui/bar/ConfidenceBar.tsx'
import DifficultyBar from '../@ui/bar/DifficultyBar.tsx'
import dayjs from 'dayjs'
import { useNavigate } from 'react-router-dom'
import { DraggableProvided } from '@hello-pangea/dnd'
import PerfectRehearsalMenuItem from '../@ui/menu/item/song/PerfectRehearsalMenuItem.tsx'
import PartialRehearsalMenuItem from '../@ui/menu/item/song/PartialRehearsalMenuItem.tsx'
import { useDeleteSongMutation } from '../../state/api/songsApi.ts'
import { useRemoveSongsFromAlbumMutation } from '../../state/api/albumsApi.ts'
import CustomIconMusicNoteEighth from '../@ui/icons/CustomIconMusicNoteEighth.tsx'
import YoutubeModal from '../@ui/modal/YoutubeModal.tsx'
import OpenLinksMenuItem from '../@ui/menu/item/song/OpenLinksMenuItem.tsx'
import AddToPlaylistMenuItem from '../@ui/menu/item/AddToPlaylistMenuItem.tsx'
import { ContextMenu } from '../@ui/menu/ContextMenu.tsx'
import useDoubleMenu from '../../hooks/useDoubleMenu.ts'
import { toast } from 'react-toastify'

interface AlbumSongCardProps {
  song: Song
  albumId: string
  isUnknownAlbum: boolean
  order: Order
  isDragging: boolean
  albumImageUrl?: string | null | undefined
  draggableProvided?: DraggableProvided
}

function AlbumSongCard({
  song,
  albumId,
  isUnknownAlbum,
  order,
  isDragging,
  albumImageUrl,
  draggableProvided
}: AlbumSongCardProps) {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const { ref: hoverRef, hovered } = useHover()
  const ref = useMergedRef(hoverRef, draggableProvided?.innerRef)

  const [removeSongsFromAlbum, { isLoading: isRemoveLoading }] = useRemoveSongsFromAlbumMutation()
  const [deleteSong, { isLoading: isDeleteLoading }] = useDeleteSongMutation()

  const { openedMenu, toggleMenu, openedContextMenu, toggleContextMenu, closeMenus } =
    useDoubleMenu()

  const isSelected = hovered || openedMenu || openedContextMenu || isDragging

  const [openedYoutube, { open: openYoutube, close: closeYoutube }] = useDisclosure(false)
  const [openedRemoveWarning, { open: openRemoveWarning, close: closeRemoveWarning }] =
    useDisclosure(false)
  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  function handleClick() {
    dispatch(openSongDrawer(song.id))
  }

  function handleViewDetails(e: MouseEvent) {
    e.stopPropagation()
    navigate(`/song/${song.id}`)
  }

  function handleOpenRemoveWarning(e: MouseEvent) {
    e.stopPropagation()
    openRemoveWarning()
  }

  function handleOpenDeleteWarning(e: MouseEvent) {
    e.stopPropagation()
    openDeleteWarning()
  }

  async function handleRemoveFromAlbum() {
    await removeSongsFromAlbum({ songIds: [song.id], id: albumId }).unwrap()
    toast.success(`${song.title} removed from album!`)
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
      <OpenLinksMenuItem song={song} openYoutube={openYoutube} />

      <Menu.Divider />
      <AddToPlaylistMenuItem ids={[song.id]} type={'song'} closeMenu={closeMenus} />
      <PartialRehearsalMenuItem songId={song.id} closeMenu={closeMenus} />
      <PerfectRehearsalMenuItem songId={song.id} closeMenu={closeMenus} />
      <Menu.Divider />

      {!isUnknownAlbum && (
        <Menu.Item leftSection={<IconCircleMinus size={14} />} onClick={handleOpenRemoveWarning}>
          Remove from Album
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
          aria-label={`song-card-${song.title}`}
          wrap={'nowrap'}
          ref={ref}
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
              boxShadow: theme.shadows.xl,
              backgroundColor: alpha(theme.colors.primary[0], 0.15)
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
                order.property === SongProperty.AlbumTrackNo ||
                order.property === SongProperty.Title
                  ? 'auto'
                  : 6
              }
            >
              <Group wrap={'nowrap'}>
                {!isUnknownAlbum && (
                  <Text fw={500} miw={25} maw={25} ta={'center'}>
                    {song.albumTrackNo}
                  </Text>
                )}
                <Avatar
                  radius={'md'}
                  src={song.imageUrl ?? albumImageUrl}
                  alt={(song.imageUrl ?? albumImageUrl) && song.title}
                  bg={'gray.5'}
                >
                  <Center c={'white'}>
                    <CustomIconMusicNoteEighth
                      aria-label={`default-icon-${song.title}`}
                      size={20}
                    />
                  </Center>
                </Avatar>

                <Text fw={500} lineClamp={1}>
                  {song.title}
                </Text>
              </Group>
            </Grid.Col>

            <Grid.Col
              span={
                order.property === SongProperty.AlbumTrackNo ||
                order.property === SongProperty.Title
                  ? 0
                  : 'auto'
              }
            >
              <Flex px={'10%'}>
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
                    aria-label={'more-menu'}
                    size={'md'}
                    variant={'grey'}
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
        title={`Remove Song From Album`}
        description={
          <Group gap={'xxs'}>
            <Text>Are you sure you want to remove</Text>
            <Text fw={600}>{song.title}</Text>
            <Text>from this album?</Text>
          </Group>
        }
        isLoading={isRemoveLoading}
        onYes={handleRemoveFromAlbum}
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

export default AlbumSongCard
