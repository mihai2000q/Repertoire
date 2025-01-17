import Album from '../../types/models/Album.ts'
import { AspectRatio, Avatar, Group, Image, Menu, Stack, Text, Title, Tooltip } from '@mantine/core'
import { IconEdit, IconInfoSquareRounded, IconTrash } from '@tabler/icons-react'
import unknownPlaceholder from '../../assets/unknown-placeholder.png'
import albumPlaceholder from '../../assets/image-placeholder-1.jpg'
import userPlaceholder from '../../assets/user-placeholder.jpg'
import dayjs from 'dayjs'
import plural from '../../utils/plural.ts'
import HeaderPanelCard from '../@ui/card/HeaderPanelCard.tsx'
import { openArtistDrawer } from '../../state/globalSlice.ts'
import { toast } from 'react-toastify'
import { useDisclosure } from '@mantine/hooks'
import { useDeleteAlbumMutation } from '../../state/albumsApi.ts'
import { useAppDispatch } from '../../state/store.ts'
import { useNavigate } from 'react-router-dom'
import AlbumInfoModal from './modal/AlbumInfoModal.tsx'
import EditAlbumHeaderModal from './modal/EditAlbumHeaderModal.tsx'
import WarningModal from '../@ui/modal/WarningModal.tsx'

interface AlbumHeaderCardProps {
  album: Album | undefined
  isUnknownAlbum: boolean
  songsTotalCount: number | undefined
}

function AlbumHeaderCard({ album, isUnknownAlbum, songsTotalCount }: AlbumHeaderCardProps) {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()

  const [deleteAlbumMutation] = useDeleteAlbumMutation()

  const [openedAlbumInfo, { open: openAlbumInfo, close: closeAlbumInfo }] = useDisclosure(false)
  const [openedEditAlbumHeader, { open: openEditAlbumHeader, close: closeEditAlbumHeader }] =
    useDisclosure(false)
  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  function handleArtistClick() {
    dispatch(openArtistDrawer(album.artist.id))
  }

  function handleDelete() {
    deleteAlbumMutation(album.id)
    navigate(`/albums`, { replace: true })
    toast.success(`${album.title} deleted!`)
  }

  return (
    <HeaderPanelCard
      onEditClick={openEditAlbumHeader}
      menuDropdown={
        <>
          <Menu.Item leftSection={<IconInfoSquareRounded size={14} />} onClick={openAlbumInfo}>
            Info
          </Menu.Item>
          <Menu.Item leftSection={<IconEdit size={14} />} onClick={openEditAlbumHeader}>
            Edit
          </Menu.Item>
          <Menu.Item leftSection={<IconTrash size={14} />} c={'red.5'} onClick={openDeleteWarning}>
            Delete
          </Menu.Item>
        </>
      }
      hideIcons={isUnknownAlbum}
    >
      <Group wrap={'nowrap'}>
        <AspectRatio>
          <Image
            h={150}
            src={isUnknownAlbum ? unknownPlaceholder : album.imageUrl}
            fallbackSrc={albumPlaceholder}
            radius={'lg'}
            alt={isUnknownAlbum ? 'unknown-album' : album.title}
            style={(theme) => ({
              boxShadow: theme.shadows.lg
            })}
          />
        </AspectRatio>
        <Stack
          gap={4}
          style={{ ...(!isUnknownAlbum && { alignSelf: 'start', paddingTop: '10px' }) }}
        >
          {!isUnknownAlbum && (
            <Text fw={500} inline>
              Album
            </Text>
          )}
          {isUnknownAlbum ? (
            <Title order={3} fw={200} fs={'italic'}>
              Unknown
            </Title>
          ) : (
            <Title order={1} fw={700} lineClamp={2}>
              {album.title}
            </Title>
          )}
          <Group gap={4}>
            {album?.artist && (
              <>
                <Group gap={'xs'}>
                  <Avatar
                    size={35}
                    src={album.artist.imageUrl ?? userPlaceholder}
                    alt={album.artist.name}
                  />
                  <Text
                    fw={600}
                    fz={'lg'}
                    sx={{
                      cursor: 'pointer',
                      '&:hover': { textDecoration: 'underline' }
                    }}
                    onClick={handleArtistClick}
                  >
                    {album.artist.name}
                  </Text>
                </Group>
                <Text c={'dimmed'}>•</Text>
              </>
            )}
            {album?.releaseDate && (
              <>
                <Tooltip
                  label={'Released on ' + dayjs(album.releaseDate).format('D MMMM YYYY')}
                  openDelay={200}
                  position={'bottom'}
                >
                  <Text fw={500} c={'dimmed'}>
                    {dayjs(album.releaseDate).format('YYYY')}
                  </Text>
                </Tooltip>
                <Text fw={500} c={'dimmed'}>
                  •
                </Text>
              </>
            )}
            <Text fw={500} c={'dimmed'}>
              {isUnknownAlbum ? songsTotalCount : album.songs.length} song
              {plural(isUnknownAlbum ? songsTotalCount : album.songs)}
            </Text>
          </Group>
        </Stack>
      </Group>

      {!isUnknownAlbum && (
        <>
          <AlbumInfoModal opened={openedAlbumInfo} onClose={closeAlbumInfo} album={album} />

          <EditAlbumHeaderModal
            album={album}
            opened={openedEditAlbumHeader}
            onClose={closeEditAlbumHeader}
          />

          <WarningModal
            opened={openedDeleteWarning}
            onClose={closeDeleteWarning}
            title={'Delete Album'}
            description={`Are you sure you want to delete this album?`}
            onYes={handleDelete}
          />
        </>
      )}
    </HeaderPanelCard>
  )
}

export default AlbumHeaderCard
