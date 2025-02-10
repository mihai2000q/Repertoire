import Album from '../../types/models/Album.ts'
import {
  AspectRatio,
  Avatar,
  Checkbox,
  Group,
  Image,
  Menu,
  Stack,
  Text,
  Title,
  Tooltip
} from '@mantine/core'
import { IconEdit, IconInfoSquareRounded, IconTrash } from '@tabler/icons-react'
import unknownPlaceholder from '../../assets/unknown-placeholder.png'
import albumPlaceholder from '../../assets/image-placeholder-1.jpg'
import userPlaceholder from '../../assets/user-placeholder.jpg'
import dayjs from 'dayjs'
import plural from '../../utils/plural.ts'
import HeaderPanelCard from '../@ui/card/HeaderPanelCard.tsx'
import { openArtistDrawer } from '../../state/slice/globalSlice.ts'
import { toast } from 'react-toastify'
import { useDisclosure } from '@mantine/hooks'
import { useDeleteAlbumMutation } from '../../state/api/albumsApi.ts'
import { useAppDispatch } from '../../state/store.ts'
import { useNavigate } from 'react-router-dom'
import AlbumInfoModal from './modal/AlbumInfoModal.tsx'
import EditAlbumHeaderModal from './modal/EditAlbumHeaderModal.tsx'
import WarningModal from '../@ui/modal/WarningModal.tsx'
import ImageModal from '../@ui/modal/ImageModal.tsx'
import { useState } from 'react'

interface AlbumHeaderCardProps {
  album: Album | undefined
  isUnknownAlbum: boolean
  songsTotalCount: number | undefined
}

function AlbumHeaderCard({ album, isUnknownAlbum, songsTotalCount }: AlbumHeaderCardProps) {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()

  const [deleteAlbumMutation, { isLoading: isDeleteLoading }] = useDeleteAlbumMutation()

  const [deleteWithSongs, setDeleteWithSongs] = useState(false)

  const [openedImage, { open: openImage, close: closeImage }] = useDisclosure(false)
  const [openedAlbumInfo, { open: openAlbumInfo, close: closeAlbumInfo }] = useDisclosure(false)
  const [openedEdit, { open: openEdit, close: closeEdit }] = useDisclosure(false)
  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  function handleArtistClick() {
    dispatch(openArtistDrawer(album.artist.id))
  }

  async function handleDelete() {
    await deleteAlbumMutation({ id: album.id, withSongs: deleteWithSongs }).unwrap()
    navigate(`/albums`, { replace: true })
    toast.success(`${album.title} deleted!`)
  }

  return (
    <HeaderPanelCard
      onEditClick={openEdit}
      menuDropdown={
        <>
          <Menu.Item leftSection={<IconInfoSquareRounded size={14} />} onClick={openAlbumInfo}>
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
      hideIcons={isUnknownAlbum}
    >
      <Group wrap={'nowrap'}>
        <AspectRatio>
          <Image
            w={150}
            src={isUnknownAlbum ? unknownPlaceholder : album.imageUrl}
            fallbackSrc={albumPlaceholder}
            radius={'lg'}
            alt={isUnknownAlbum ? 'unknown-album' : album.title}
            sx={(theme) => ({
              boxShadow: theme.shadows.lg,
              ...(!isUnknownAlbum && album.imageUrl && { cursor: 'pointer' })
            })}
            onClick={!isUnknownAlbum && album.imageUrl ? openImage : undefined}
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
          <Group gap={4} wrap={'nowrap'}>
            {album?.artist && (
              <>
                <Group gap={'xs'} wrap={'nowrap'}>
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
                    lineClamp={1}
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
          <ImageModal
            opened={openedImage}
            onClose={closeImage}
            title={album.title}
            image={album.imageUrl}
          />

          <AlbumInfoModal opened={openedAlbumInfo} onClose={closeAlbumInfo} album={album} />

          <EditAlbumHeaderModal album={album} opened={openedEdit} onClose={closeEdit} />

          <WarningModal
            opened={openedDeleteWarning}
            onClose={closeDeleteWarning}
            title={'Delete Album'}
            description={
              <Stack gap={5}>
                <Text fw={500}>Are you sure you want to delete this album?</Text>
                <Checkbox
                  checked={deleteWithSongs}
                  onChange={(event) => setDeleteWithSongs(event.currentTarget.checked)}
                  label={'Delete all associated songs'}
                  c={'dimmed'}
                  styles={{ label: { paddingLeft: 8 } }}
                />
              </Stack>
            }
            onYes={handleDelete}
            isLoading={isDeleteLoading}
          />
        </>
      )}
    </HeaderPanelCard>
  )
}

export default AlbumHeaderCard
