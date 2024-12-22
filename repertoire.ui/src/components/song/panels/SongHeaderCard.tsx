import Song from '../../../types/models/Song.ts'
import {
  AspectRatio,
  Avatar,
  Group,
  HoverCard,
  Image,
  Menu,
  Stack,
  Text,
  Title,
  Tooltip
} from '@mantine/core'
import { IconEdit, IconInfoSquareRounded, IconTrash } from '@tabler/icons-react'
import songPlaceholder from '../../../assets/image-placeholder-1.jpg'
import userPlaceholder from '../../../assets/user-placeholder.jpg'
import dayjs from 'dayjs'
import HeaderPanelCard from '../../@ui/card/HeaderPanelCard.tsx'
import EditSongHeaderModal from '../modal/EditSongHeaderModal.tsx'
import { useDisclosure } from '@mantine/hooks'
import { openAlbumDrawer, openArtistDrawer } from '../../../state/globalSlice.ts'
import { toast } from 'react-toastify'
import { useDeleteSongMutation } from '../../../state/songsApi.ts'
import { useAppDispatch } from '../../../state/store.ts'
import { useNavigate } from 'react-router-dom'
import SongInfoModal from '../modal/SongInfoModal.tsx'
import WarningModal from '../../@ui/modal/WarningModal.tsx'

interface SongHeaderCardProps {
  song: Song
}

function SongHeaderCard({ song }: SongHeaderCardProps) {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()

  const [deleteSongMutation] = useDeleteSongMutation()

  const [openedEdit, { open: openEdit, close: closeEdit }] = useDisclosure(false)
  const [openedInfo, { open: openInfo, close: closeInfo }] = useDisclosure(false)
  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  function handleAlbumClick() {
    dispatch(openAlbumDrawer(song.album.id))
  }

  function handleArtistClick() {
    dispatch(openArtistDrawer(song.artist.id))
  }

  function handleDelete() {
    deleteSongMutation(song.id)
    navigate(`/songs`, { replace: true })
    toast.success(`${song.title} deleted!`)
  }

  return (
    <HeaderPanelCard
      onEditClick={openEdit}
      menuDropdown={
        <>
          <Menu.Item leftSection={<IconInfoSquareRounded size={14} />} onClick={openInfo}>
            Info
          </Menu.Item>
          <Menu.Item leftSection={<IconEdit size={14} />} onClick={openEdit}>
            Edit
          </Menu.Item>
          <Menu.Item leftSection={<IconTrash size={14} />} c={'red.5'} onClick={openDeleteWarning}>
            Delete
          </Menu.Item>
        </>
      }
    >
      <Group wrap={'nowrap'}>
        <AspectRatio>
          <Image
            h={150}
            src={song.imageUrl}
            fallbackSrc={songPlaceholder}
            radius={'lg'}
            style={(theme) => ({
              boxShadow: theme.shadows.lg
            })}
          />
        </AspectRatio>
        <Stack gap={4} style={{ alignSelf: 'start' }} pt={'xs'}>
          <Text fw={500} inline>
            Song
          </Text>
          <Title order={1} fw={700} lineClamp={2}>
            {song.title}
          </Title>
          <Group gap={4}>
            {song.artist && (
              <Group gap={'xs'}>
                <Avatar size={35} src={song.artist.imageUrl ?? userPlaceholder} />
                <Text
                  fw={700}
                  fz={'lg'}
                  sx={{
                    cursor: 'pointer',
                    '&:hover': { textDecoration: 'underline' }
                  }}
                  inline
                  onClick={handleArtistClick}
                >
                  {song.artist.name}
                </Text>
              </Group>
            )}

            {song.album && (
              <Group gap={0}>
                {song.artist && (
                  <Text fw={500} c={'dimmed'} inline pr={4}>
                    on
                  </Text>
                )}
                <HoverCard shadow={'lg'} withArrow>
                  <HoverCard.Target>
                    <Text
                      fw={600}
                      inline
                      c={'dark'}
                      sx={{
                        cursor: 'pointer',
                        '&:hover': { textDecoration: 'underline' }
                      }}
                      onClick={handleAlbumClick}
                    >
                      {song.album.title}
                    </Text>
                  </HoverCard.Target>
                  <HoverCard.Dropdown maw={300}>
                    <Group align={'center'} gap={'xs'} wrap={'nowrap'}>
                      <Avatar
                        size={45}
                        radius={'md'}
                        src={song.album.imageUrl ?? songPlaceholder}
                      />
                      <Stack gap={2}>
                        <Text fw={500} fz={'xs'} inline>
                          Album
                        </Text>
                        <Text fw={600} fz={'md'} inline lineClamp={2}>
                          {song.album.title}
                        </Text>
                        {song.album.releaseDate && (
                          <Text fw={500} c={'dimmed'} fz={'sm'} inline>
                            {dayjs(song.album.releaseDate).format('DD MMM YYYY')}
                          </Text>
                        )}
                      </Stack>
                    </Group>
                  </HoverCard.Dropdown>
                </HoverCard>
              </Group>
            )}

            {song.releaseDate && (
              <>
                {(song.album || song.artist) && (
                  <Text c={'dimmed'} inline>
                    •
                  </Text>
                )}
                <Tooltip
                  label={'Released on ' + dayjs(song.releaseDate).format('DD MMMM YYYY')}
                  openDelay={200}
                  position={'bottom'}
                >
                  <Text fw={500} c={'dimmed'} inline>
                    {dayjs(song.releaseDate).format('YYYY')}
                  </Text>
                </Tooltip>
              </>
            )}
          </Group>
        </Stack>
      </Group>

      <SongInfoModal opened={openedInfo} onClose={closeInfo} song={song} />
      <EditSongHeaderModal song={song} opened={openedEdit} onClose={closeEdit} />

      <WarningModal
        opened={openedDeleteWarning}
        onClose={closeDeleteWarning}
        title={'Delete Song'}
        description={`Are you sure you want to delete this song?`}
        onYes={handleDelete}
      />
    </HeaderPanelCard>
  )
}

export default SongHeaderCard
