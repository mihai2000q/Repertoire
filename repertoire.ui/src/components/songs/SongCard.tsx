import Song from '../../types/models/Song'
import { Anchor, Avatar, Box, Card, Center, Group, Menu, Stack, Text, Tooltip } from '@mantine/core'
import { useNavigate } from 'react-router-dom'
import { useAppDispatch } from '../../state/store.ts'
import { openArtistDrawer } from '../../state/slice/globalSlice.ts'
import { MouseEvent, ReactElement } from 'react'
import {
  IconBoltFilled,
  IconBombFilled,
  IconBrandYoutubeFilled,
  IconGuitarPickFilled,
  IconMicrophoneFilled,
  IconStarFilled,
  IconTrash
} from '@tabler/icons-react'
import useDifficultyInfo from '../../hooks/useDifficultyInfo.ts'
import { toast } from 'react-toastify'
import { useDeleteSongMutation } from '../../state/api/songsApi.ts'
import useContextMenu from '../../hooks/useContextMenu.ts'
import { useDisclosure } from '@mantine/hooks'
import WarningModal from '../@ui/modal/WarningModal.tsx'
import CustomIconGuitarHead from '../@ui/icons/CustomIconGuitarHead.tsx'
import CustomIconLightningTrio from '../@ui/icons/CustomIconLightningTrio.tsx'
import YoutubeModal from '../@ui/modal/YoutubeModal.tsx'
import PerfectRehearsalMenuItem from '../@ui/menu/item/PerfectRehearsalMenuItem.tsx'
import PartialRehearsalMenuItem from '../@ui/menu/item/PartialRehearsalMenuItem.tsx'
import CustomIconMusicNoteEighth from '../@ui/icons/CustomIconMusicNoteEighth.tsx'

const iconSize = 18
const LocalAnchor = ({ link, children }: { link: string; children: ReactElement }) => (
  <Anchor
    underline={'never'}
    href={link}
    target="_blank"
    rel="noreferrer"
    c={'inherit'}
    onClick={(e) => e.stopPropagation()}
  >
    {children}
  </Anchor>
)

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

  const [openedMenu, menuDropdownProps, { openMenu, closeMenu }] = useContextMenu()

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

  async function handleDelete() {
    await deleteSongMutation(song.id).unwrap()
    toast.success(`${song.title} deleted!`)
  }

  return (
    <Menu shadow={'lg'} opened={openedMenu} onClose={closeMenu}>
      <Menu.Target>
        <Card
          aria-label={`song-card-${song.title}`}
          p={0}
          radius={'10%'}
          onClick={handleClick}
          onContextMenu={openMenu}
          sx={(theme) => ({
            cursor: 'pointer',
            transition: '0.3s',
            boxShadow: theme.shadows.lg,
            '&:hover': {
              boxShadow: theme.shadows.xxl,
              transform: 'scale(1.1)'
            }
          })}
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
                  {song?.riffs > 1 && (
                    <LocalTooltip label={`This song has ${song.riffs} riffs`}>
                      <IconBombFilled size={iconSize} aria-label={'riffs-icon'} />
                    </LocalTooltip>
                  )}
                  {song?.solos > 0 && (
                    <LocalTooltip
                      label={
                        song.solos === 1
                          ? 'This song has a solo'
                          : `This song has ${song.solos} solos`
                      }
                    >
                      <Center c={song.solos === 1 ? 'yellow.4' : 'yellow.5'}>
                        {song.solos > 1 ? (
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
                    <LocalAnchor link={song.songsterrLink}>
                      <LocalTooltip label={'Open Songsterr'}>
                        <Center c={'blue.7'}>
                          <IconGuitarPickFilled size={iconSize} aria-label={'songsterr-icon'} />
                        </Center>
                      </LocalTooltip>
                    </LocalAnchor>
                  )}
                  {song.youtubeLink && (
                    <LocalTooltip label={'Open Youtube'}>
                      <Center c={'red.7'} onClick={handleOpenYoutube}>
                        <IconBrandYoutubeFilled size={iconSize} aria-label={'youtube-icon'} />
                      </Center>
                    </LocalTooltip>
                  )}
                </Tooltip.Group>
              </Group>
            </Stack>
          </Stack>
        </Card>
      </Menu.Target>

      <Menu.Dropdown {...menuDropdownProps}>
        <PartialRehearsalMenuItem songId={song.id} />
        <PerfectRehearsalMenuItem songId={song.id} />
        <Menu.Item c={'red'} leftSection={<IconTrash size={14} />} onClick={openDeleteWarning}>
          Delete
        </Menu.Item>
      </Menu.Dropdown>

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
    </Menu>
  )
}

export default SongCard
