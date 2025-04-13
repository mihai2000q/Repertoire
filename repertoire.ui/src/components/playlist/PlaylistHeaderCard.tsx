import { Avatar, Center, Group, Menu, Stack, Text, Title } from '@mantine/core'
import { IconEdit, IconInfoSquareRounded, IconPlaylist, IconTrash } from '@tabler/icons-react'
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
import titleFontSize from '../../utils/titleFontSize.ts'

interface PlaylistHeaderCardProps {
  playlist: Playlist
}

function PlaylistHeaderCard({ playlist }: PlaylistHeaderCardProps) {
  const navigate = useNavigate()

  const [deletePlaylistMutation, { isLoading: isDeleteLoading }] = useDeletePlaylistMutation()

  const [openedImage, { open: openImage, close: closeImage }] = useDisclosure(false)
  const [openedInfo, { open: openInfo, close: closeInfo }] = useDisclosure(false)
  const [openedEdit, { open: openEdit, close: closeEdit }] = useDisclosure(false)
  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  async function handleDelete() {
    await deletePlaylistMutation(playlist.id).unwrap()
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
        <Avatar
          src={playlist.imageUrl}
          alt={playlist.imageUrl && playlist.title}
          size={'max(12vw, 150px)'}
          radius={'10%'}
          bg={'gray.5'}
          style={(theme) => ({
            boxShadow: theme.shadows.lg,
            ...(playlist.imageUrl && { cursor: 'pointer' })
          })}
          onClick={playlist.imageUrl && openImage}
        >
          <Center c={'white'}>
            <IconPlaylist
              aria-label={`default-icon-${playlist.title}`}
              size={'100%'}
              style={{ padding: '28%' }}
            />
          </Center>
        </Avatar>

        <Stack gap={'xxs'} pt={'md'}>
          <Text fw={500} inline>
            Playlist
          </Text>

          <Title order={1} fw={700} lineClamp={2} fz={titleFontSize(playlist.title)}>
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
        isLoading={isDeleteLoading}
      />
    </HeaderPanelCard>
  )
}

export default PlaylistHeaderCard
