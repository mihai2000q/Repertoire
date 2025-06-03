import Song from '../../../types/models/Song.ts'
import { Avatar, Center, Group, HoverCard, Menu, Stack, Text, Title, Tooltip } from '@mantine/core'
import { IconEdit, IconInfoSquareRounded, IconTrash } from '@tabler/icons-react'
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
import PerfectRehearsalMenuItem from '../../@ui/menu/item/song/PerfectRehearsalMenuItem.tsx'
import titleFontSize from '../../../utils/style/titleFontSize.ts'
import PartialRehearsalMenuItem from '../../@ui/menu/item/song/PartialRehearsalMenuItem.tsx'
import CustomIconMusicNoteEighth from '../../@ui/icons/CustomIconMusicNoteEighth.tsx'
import CustomIconAlbumVinyl from '../../@ui/icons/CustomIconAlbumVinyl.tsx'
import CustomIconUserAlt from '../../@ui/icons/CustomIconUserAlt.tsx'

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

  const [openedMenu, { open: openMenu, close: closeMenu }] = useDisclosure(false)

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
      menuOpened={openedMenu}
      openMenu={openMenu}
      closeMenu={closeMenu}
      menuDropdown={
        <>
          <Menu.Item leftSection={<IconInfoSquareRounded size={14} />} onClick={openInfo}>
            Info
          </Menu.Item>
          <Menu.Item leftSection={<IconEdit size={14} />} onClick={openEdit}>
            Edit
          </Menu.Item>
          <PartialRehearsalMenuItem songId={song.id} />
          <PerfectRehearsalMenuItem songId={song.id} />
          <Menu.Item leftSection={<IconTrash size={14} />} c={'red.5'} onClick={openDeleteWarning}>
            Delete
          </Menu.Item>
        </>
      }
    >
      <Group wrap={'nowrap'}>
        <Avatar
          radius={'10%'}
          src={song.imageUrl ?? song.album?.imageUrl}
          alt={(song.imageUrl ?? song.album?.imageUrl) && song.title}
          size={'max(12vw, 150px)'}
          bg={'gray.5'}
          style={(theme) => ({
            aspectRatio: 1,
            boxShadow: theme.shadows.lg,
            ...((song.imageUrl || song.album?.imageUrl) && { cursor: 'pointer' })
          })}
          onClick={(song.imageUrl || song.album?.imageUrl) && openImage}
        >
          <Center c={'white'}>
            <CustomIconMusicNoteEighth
              aria-label={`default-icon-${song.title}`}
              size={'100%'}
              style={{ padding: '26%' }}
            />
          </Center>
        </Avatar>

        <Stack gap={'xxs'}>
          <Text fw={500} inline>
            Song
          </Text>
          <Title order={1} fw={700} lineClamp={2} fz={titleFontSize(song.title)}>
            {song.title}
          </Title>
          <Group gap={'xxs'} wrap={'nowrap'}>
            {song.artist && (
              <Group gap={'xs'} wrap={'nowrap'}>
                <Avatar
                  size={35}
                  src={song.artist.imageUrl}
                  alt={song.artist.imageUrl && song.artist.name}
                  style={(theme) => ({ boxShadow: theme.shadows.sm })}
                  bg={'gray.0'}
                >
                  <Center c={'gray.7'}>
                    <CustomIconUserAlt aria-label={`default-icon-${song.artist.name}`} size={15} />
                  </Center>
                </Avatar>
                <Text
                  fw={700}
                  fz={'lg'}
                  sx={{
                    cursor: 'pointer',
                    '&:hover': { textDecoration: 'underline' }
                  }}
                  lh={'xxs'}
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
                <HoverCard>
                  <HoverCard.Target>
                    <Text
                      fw={600}
                      lh={'xxs'}
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
                        radius={'md'}
                        size={45}
                        src={song.album.imageUrl}
                        alt={song.album.imageUrl && song.album.title}
                        bg={'gray.5'}
                      >
                        <Center c={'white'}>
                          <CustomIconAlbumVinyl
                            aria-label={`default-icon-${song.album.imageUrl}`}
                            size={18}
                          />
                        </Center>
                      </Avatar>
                      <Stack gap={2}>
                        <Text fw={500} fz={'xs'} inline>
                          Album
                        </Text>
                        <Text fw={600} fz={'md'} lh={'xxs'} lineClamp={2}>
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
