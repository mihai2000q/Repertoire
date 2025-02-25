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
import albumPlaceholder from '../../../assets/image-placeholder-1.jpg'
import userPlaceholder from '../../../assets/user-placeholder.jpg'
import dayjs from 'dayjs'
import HeaderPanelCard from '../../@ui/card/HeaderPanelCard.tsx'
import EditSongHeaderModal from '../modal/EditSongHeaderModal.tsx'
import { useDisclosure } from '@mantine/hooks'
import { openAlbumDrawer, openArtistDrawer } from '../../../state/slice/globalSlice.ts'
import { toast } from 'react-toastify'
import { useDeleteSongMutation } from '../../../state/api/songsApi.ts'
import { useAppDispatch } from '../../../state/store.ts'
import { useNavigate } from 'react-router-dom'
import SongInfoModal from '../modal/SongInfoModal.tsx'
import WarningModal from '../../@ui/modal/WarningModal.tsx'
import ImageModal from '../../@ui/modal/ImageModal.tsx'
import PerfectRehearsalMenuItem from '../../@ui/menu/item/PerfectRehearsalMenuItem.tsx'
import titleFontSize from "../../../utils/titleFontSize.ts";

interface SongHeaderCardProps {
  song: Song
}

function SongHeaderCard({ song }: SongHeaderCardProps) {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()

  const [deleteSongMutation, { isLoading: isDeleteLoading }] = useDeleteSongMutation()

  const [openedImage, { open: openImage, close: closeImage }] = useDisclosure(false)
  const [openedInfo, { open: openInfo, close: closeInfo }] = useDisclosure(false)
  const [openedEdit, { open: openEdit, close: closeEdit }] = useDisclosure(false)
  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  function handleAlbumClick() {
    dispatch(openAlbumDrawer(song.album.id))
  }

  function handleArtistClick() {
    dispatch(openArtistDrawer(song.artist.id))
  }

  async function handleDelete() {
    await deleteSongMutation(song.id).unwrap()
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
          <PerfectRehearsalMenuItem songId={song.id} />
          <Menu.Item leftSection={<IconTrash size={14} />} c={'red.5'} onClick={openDeleteWarning}>
            Delete
          </Menu.Item>
        </>
      }
    >
      <Group wrap={'nowrap'}>
        <AspectRatio>
          <Image
            w={'max(12vw, 150px)'}
            src={song.imageUrl ?? song.album?.imageUrl}
            fallbackSrc={songPlaceholder}
            alt={song.title}
            radius={'lg'}
            sx={(theme) => ({
              boxShadow: theme.shadows.lg,
              ...((song.imageUrl || song.album?.imageUrl) && { cursor: 'pointer' })
            })}
            onClick={(song.imageUrl || song.album?.imageUrl) && openImage}
          />
        </AspectRatio>
        <Stack gap={'xxs'}>
          <Text fw={500} inline>
            Song
          </Text>
          <Title
            order={1}
            fw={700}
            lineClamp={2}
            fz={titleFontSize(song.title)}
          >
            {song.title}
          </Title>
          <Group gap={'xxs'} wrap={'nowrap'}>
            {song.artist && (
              <Group gap={'xs'} wrap={'nowrap'}>
                <Avatar
                  size={35}
                  src={song.artist.imageUrl ?? userPlaceholder}
                  alt={song.artist.name}
                />
                <Text
                  fw={700}
                  fz={'lg'}
                  sx={{
                    cursor: 'pointer',
                    '&:hover': { textDecoration: 'underline' }
                  }}
                  inline
                  onClick={handleArtistClick}
                  lineClamp={1}
                >
                  {song.artist.name}
                </Text>
              </Group>
            )}

            {song.album && (
              <Group gap={0} wrap={'nowrap'}>
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
                      lineClamp={1}
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
                    <Group gap={'xs'} wrap={'nowrap'}>
                      <Avatar
                        size={45}
                        radius={'md'}
                        src={song.album.imageUrl ?? albumPlaceholder}
                        alt={song.album.title}
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
                            {dayjs(song.album.releaseDate).format('D MMM YYYY')}
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
                  label={'Released on ' + dayjs(song.releaseDate).format('D MMMM YYYY')}
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

      <ImageModal
        opened={openedImage}
        onClose={closeImage}
        title={song.title}
        image={song.imageUrl || song.album?.imageUrl}
      />

      <SongInfoModal opened={openedInfo} onClose={closeInfo} song={song} />
      <EditSongHeaderModal song={song} opened={openedEdit} onClose={closeEdit} />
      <WarningModal
        opened={openedDeleteWarning}
        onClose={closeDeleteWarning}
        title={'Delete Song'}
        description={`Are you sure you want to delete this song?`}
        onYes={handleDelete}
        isLoading={isDeleteLoading}
      />
    </HeaderPanelCard>
  )
}

export default SongHeaderCard
