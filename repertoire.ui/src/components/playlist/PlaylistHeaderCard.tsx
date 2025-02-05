import { AspectRatio, Group, Image, Menu, Stack, Text, Title } from '@mantine/core'
import { IconEdit, IconInfoSquareRounded, IconTrash } from '@tabler/icons-react'
import playlistPlaceholder from '../../assets/image-placeholder-1.jpg'
import plural from '../../utils/plural.ts'
import HeaderPanelCard from '../@ui/card/HeaderPanelCard.tsx'
import Playlist from '../../types/models/Playlist.ts'
import { useDisclosure } from '@mantine/hooks'
import { toast } from 'react-toastify'
import { useDeletePlaylistMutation } from '../../state/api/playlistsApi.ts'
import { useNavigate } from 'react-router-dom'
import EditPlaylistHeaderModal from './modal/EditPlaylistHeaderModal.tsx'
import PlaylistInfoModal from './modal/PlaylistInfoModal.tsx'
import WarningModal from '../@ui/modal/WarningModal.tsx'
import ImageModal from '../@ui/modal/ImageModal.tsx'

interface PlaylistHeaderCardProps {
  playlist: Playlist
}

function PlaylistHeaderCard({ playlist }: PlaylistHeaderCardProps) {
  const navigate = useNavigate()

  const [deletePlaylistMutation] = useDeletePlaylistMutation()

  const [openedImage, { open: openImage, close: closeImage }] = useDisclosure(false)
  const [openedInfo, { open: openInfo, close: closeInfo }] = useDisclosure(false)
  const [openedEdit, { open: openEdit, close: closeEdit }] = useDisclosure(false)
  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  function handleDelete() {
    deletePlaylistMutation(playlist.id)
    navigate(`/playlists`, { replace: true })
    toast.success(`${playlist.title} deleted!`)
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
            w={150}
            src={playlist.imageUrl}
            alt={playlist.title}
            fallbackSrc={playlistPlaceholder}
            radius={'lg'}
            sx={(theme) => ({
              boxShadow: theme.shadows.lg,
              ...(playlist.imageUrl && { cursor: 'pointer' })
            })}
            onClick={playlist.imageUrl && openImage}
          />
        </AspectRatio>
        <Stack gap={4} pt={'md'} style={{ alignSelf: 'start' }}>
          <Text fw={500} inline>
            Playlist
          </Text>

          <Title order={1} fw={700} lineClamp={2}>
            {playlist.title}
          </Title>

          <Text fw={500} fz={'sm'} c={'dimmed'} inline>
            {playlist.songs.length} song{plural(playlist.songs)}
          </Text>

          <Text fz={'sm'} c={'dimmed'} lineClamp={3}>
            {playlist.description}
          </Text>
        </Stack>
      </Group>

      <ImageModal
        opened={openedImage}
        onClose={closeImage}
        title={playlist.title}
        image={playlist.imageUrl}
      />

      <PlaylistInfoModal playlist={playlist} opened={openedInfo} onClose={closeInfo} />
      <EditPlaylistHeaderModal playlist={playlist} opened={openedEdit} onClose={closeEdit} />
      <WarningModal
        opened={openedDeleteWarning}
        onClose={closeDeleteWarning}
        title={'Delete Playlist'}
        description={`Are you sure you want to delete this playlist?`}
        onYes={handleDelete}
      />
    </HeaderPanelCard>
  )
}

export default PlaylistHeaderCard
