import Song from '../../types/models/Song'
import { Anchor, Avatar, Box, Card, Center, Group, Stack, Text, Tooltip } from '@mantine/core'
import { useNavigate } from 'react-router-dom'
import { useAppDispatch } from '../../state/store.ts'
import { openArtistDrawer, openSongDrawer } from '../../state/slice/globalSlice.ts'
import { MouseEvent, ReactElement } from 'react'
import {
  IconBoltFilled,
  IconBombFilled,
  IconBrandYoutubeFilled,
  IconDisc,
  IconGuitarPickFilled,
  IconLayoutSidebarLeftExpand,
  IconMicrophoneFilled,
  IconStarFilled,
  IconTrash,
  IconUser
} from '@tabler/icons-react'
import useDifficultyInfo from '../../hooks/useDifficultyInfo.ts'
import { toast } from 'react-toastify'
import { useDeleteSongMutation } from '../../state/api/songsApi.ts'
import { useDisclosure } from '@mantine/hooks'
import WarningModal from '../@ui/modal/WarningModal.tsx'
import CustomIconGuitarHead from '../@ui/icons/CustomIconGuitarHead.tsx'
import CustomIconLightningTrio from '../@ui/icons/CustomIconLightningTrio.tsx'
import YoutubeModal from '../@ui/modal/YoutubeModal.tsx'
import PerfectRehearsalMenuItem from '../@ui/menu/item/PerfectRehearsalMenuItem.tsx'
import PartialRehearsalMenuItem from '../@ui/menu/item/song/PartialRehearsalMenuItem.tsx'
import CustomIconMusicNoteEighth from '../@ui/icons/CustomIconMusicNoteEighth.tsx'
import AddToPlaylistMenuItem from '../@ui/menu/item/AddToPlaylistMenuItem.tsx'
import { ContextMenu } from '../@ui/menu/ContextMenu.tsx'

const iconSize = 18

const LocalTooltip = ({ label, children }: { label: string; children: ReactElement }) => (
  <Tooltip label={label} position="bottom">
    {children}
  </Tooltip>
)

interface SongCardProps {
  song: Song
}

function SongCard({ song }: SongCardProps) {
  const navigate = useNavigate()
  const dispatch = useAppDispatch()

  const [deleteSongMutation, { isLoading: isDeleteLoading }] = useDeleteSongMutation()

  const { color: difficultyColor } = useDifficultyInfo(song?.difficulty)

  const [openedMenu, { open: openMenu, close: closeMenu }] = useDisclosure(false)

  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)
  const [openedYoutube, { open: openYoutube, close: closeYoutube }] = useDisclosure(false)

  function handleClick() {
    navigate(`/song/${song.id}`)
  }

  function handleArtistClick(e: MouseEvent<HTMLElement>) {
    e.stopPropagation()
    dispatch(openArtistDrawer(song.artist.id))
  }

  function handleOpenYoutube(e: MouseEvent<HTMLElement>) {
    e.stopPropagation()
    openYoutube()
  }

  function handleOpenDrawer() {
    dispatch(openSongDrawer(song.id))
  }

  function handleViewArtist() {
    navigate(`/artist/${song.artist.id}`)
  }

  function handleViewAlbum() {
    navigate(`/album/${song.album.id}`)
  }

  async function handleDelete() {
    await deleteSongMutation(song.id).unwrap()
    toast.success(`${song.title} deleted!`)
  }

  return (
    <ContextMenu opened={openedMenu} onClose={closeMenu} onOpen={openMenu}>
      <ContextMenu.Target>
        <Card
          aria-label={`song-card-${song.title}`}
          p={0}
          radius={'10%'}
          sx={(theme) => ({
            cursor: 'pointer',
            transition: '0.3s',
            boxShadow: theme.shadows.lg,
            '&:hover': {
              boxShadow: theme.shadows.xxl,
              transform: 'scale(1.1)'
            },
            ...(openedMenu && {
              boxShadow: theme.shadows.xxl,
              transform: 'scale(1.1)'
            })
          })}
          onClick={handleClick}
        >
          <Stack gap={0}>
            <Avatar
              radius={'10%'}
              src={song.imageUrl ?? song.album?.imageUrl}
              alt={(song.imageUrl ?? song.album?.imageUrl) && song.title}
              w={'100%'}
              h={'unset'}
              bg={'gray.5'}
              style={(theme) => ({
                aspectRatio: 8 / 7,
                boxShadow: theme.shadows.sm
              })}
            >
              <Center c={'white'}>
                <CustomIconMusicNoteEighth
                  aria-label={`default-icon-${song.title}`}
                  size={'100%'}
                  style={{ padding: '30%' }}
                />
              </Center>
            </Avatar>

            <Stack gap={0} px={'sm'} pt={'xs'} pb={6} align={'start'}>
              <Text fw={600} lineClamp={2} lh={'xxs'} mb={1}>
                {song.title}
              </Text>
              <Box pb={1}>
                {song.artist ? (
                  <Text
                    fz={'sm'}
                    c="dimmed"
                    lineClamp={1}
                    onClick={handleArtistClick}
                    sx={{ '&:hover': { textDecoration: 'underline' } }}
                  >
                    {song.artist?.name}
                  </Text>
                ) : (
                  <Text fz={'sm'} c="dimmed" fs={'oblique'}>
                    Unknown
                  </Text>
                )}
              </Box>
              <Group wrap={'nowrap'} c={'primary.9'} gap={'xxs'} align={'end'} pb={1}>
                <Tooltip.Group openDelay={200}>
                  {song.isRecorded && (
                    <LocalTooltip label={'This song is recorded'}>
                      <IconMicrophoneFilled size={iconSize - 2} aria-label={'recorded-icon'} />
                    </LocalTooltip>
                  )}
                  {song.guitarTuning && (
                    <LocalTooltip label={`This song is tuned in ${song.guitarTuning.name}`}>
                      <CustomIconGuitarHead size={iconSize} aria-label={'guitar-tuning-icon'} />
                    </LocalTooltip>
                  )}
                  {song.riffsCount > 1 && (
                    <LocalTooltip label={`This song has ${song.riffsCount} riffs`}>
                      <IconBombFilled size={iconSize} aria-label={'riffs-icon'} />
                    </LocalTooltip>
                  )}
                  {song.solosCount > 0 && (
                    <LocalTooltip
                      label={
                        song.solosCount === 1
                          ? 'This song has a solo'
                          : `This song has ${song.solosCount} solos`
                      }
                    >
                      <Center c={song.solosCount === 1 ? 'yellow.4' : 'yellow.5'}>
                        {song.solosCount > 1 ? (
                          <CustomIconLightningTrio size={iconSize} aria-label={'solos-icon'} />
                        ) : (
                          <IconBoltFilled size={iconSize} aria-label={'solo-icon'} />
                        )}
                      </Center>
                    </LocalTooltip>
                  )}
                  {song.difficulty && (
                    <LocalTooltip label={`This song is ${song.difficulty}`}>
                      <Center c={difficultyColor}>
                        <IconStarFilled size={iconSize} aria-label={'difficulty-icon'} />
                      </Center>
                    </LocalTooltip>
                  )}
                  {song.songsterrLink && (
                    <Anchor
                      underline={'never'}
                      aria-label={'songsterr'}
                      href={song.songsterrLink}
                      target="_blank"
                      rel="noreferrer"
                      c={'inherit'}
                      onClick={(e) => e.stopPropagation()}
                    >
                      <LocalTooltip label={'Open Songsterr'}>
                        <Center c={'blue.7'}>
                          <IconGuitarPickFilled
                            role={'button'}
                            size={iconSize}
                            aria-label={'songsterr'}
                          />
                        </Center>
                      </LocalTooltip>
                    </Anchor>
                  )}
                  {song.youtubeLink && (
                    <LocalTooltip label={'Open Youtube'}>
                      <Center c={'red.7'} onClick={handleOpenYoutube}>
                        <IconBrandYoutubeFilled
                          role={'button'}
                          size={iconSize}
                          aria-label={'youtube'}
                        />
                      </Center>
                    </LocalTooltip>
                  )}
                </Tooltip.Group>
              </Group>
            </Stack>
          </Stack>
        </Card>
      </ContextMenu.Target>

      <ContextMenu.Dropdown>
        <ContextMenu.Item
          leftSection={<IconLayoutSidebarLeftExpand size={14} />}
          onClick={handleOpenDrawer}
        >
          Open Drawer
        </ContextMenu.Item>
        <ContextMenu.Item
          leftSection={<IconUser size={14} />}
          disabled={!song.artist}
          onClick={handleViewArtist}
        >
          View Artist
        </ContextMenu.Item>
        <ContextMenu.Item
          leftSection={<IconDisc size={14} />}
          disabled={!song.album}
          onClick={handleViewAlbum}
        >
          View Album
        </ContextMenu.Item>

        <ContextMenu.Divider />
        <AddToPlaylistMenuItem ids={[song.id]} type={'song'} closeMenu={closeMenu} />
        <PartialRehearsalMenuItem songId={song.id} closeMenu={closeMenu} />
        <PerfectRehearsalMenuItem id={song.id} closeMenu={closeMenu} type={'song'} />
        <ContextMenu.Divider />

        <ContextMenu.Item
          c={'red'}
          leftSection={<IconTrash size={14} />}
          onClick={openDeleteWarning}
        >
          Delete
        </ContextMenu.Item>
      </ContextMenu.Dropdown>

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
        onYes={handleDelete}
        isLoading={isDeleteLoading}
      />
      <YoutubeModal
        title={song.title}
        link={song.youtubeLink}
        opened={openedYoutube}
        onClose={closeYoutube}
      />
    </ContextMenu>
  )
}

export default SongCard
