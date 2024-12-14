import { AspectRatio, Group, Image, Menu, Stack, Text, Title } from '@mantine/core'
import { IconEdit, IconInfoSquareRounded, IconTrash } from '@tabler/icons-react'
import playlistPlaceholder from '../../assets/image-placeholder-1.jpg'
import plural from '../../utils/plural.ts'
import HeaderPanelCard from '../card/HeaderPanelCard.tsx'
import Playlist from '../../types/models/Playlist.ts'
import { useDisclosure } from '@mantine/hooks'
import { toast } from 'react-toastify'
import { useDeletePlaylistMutation } from '../../state/playlistsApi.ts'
import { useNavigate } from 'react-router-dom'
import EditPlaylistHeaderModal from './modal/EditPlaylistHeaderModal.tsx'
import PlaylistInfoModal from './modal/PlaylistInfoModal.tsx'

interface PlaylistHeaderCardProps {
  playlist: Playlist
}

function PlaylistHeaderCard({ playlist }: PlaylistHeaderCardProps) {
  const navigate = useNavigate()

  const [deletePlaylistMutation] = useDeletePlaylistMutation()

  const [openedInfo, { open: openInfo, close: closeInfo }] = useDisclosure(false)
  const [openedEdit, { open: openEdit, close: closeEdit }] = useDisclosure(false)

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
            Edit
          </Menu.Item>
          <Menu.Item leftSection={<IconEdit size={14} />} onClick={openEdit}>
            Edit
          </Menu.Item>
          <Menu.Item leftSection={<IconTrash size={14} />} c={'red.5'} onClick={handleDelete}>
            Delete
          </Menu.Item>
        </>
      }
    >
      <Group>
        <AspectRatio>
          <Image
            h={150}
            src={playlist.imageUrl}
            fallbackSrc={playlistPlaceholder}
            radius={'lg'}
            style={(theme) => ({
              boxShadow: theme.shadows.lg
            })}
          />
        </AspectRatio>
        <Stack gap={4} style={{ alignSelf: 'start', paddingTop: '10px' }}>
          <Text fw={500} inline>
            Playlist
          </Text>

          <Title order={1} fw={700}>
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

      <PlaylistInfoModal playlist={playlist} opened={openedInfo} onClose={closeInfo} />
      <EditPlaylistHeaderModal playlist={playlist} opened={openedEdit} onClose={closeEdit} />
    </HeaderPanelCard>
  )
}

export default PlaylistHeaderCard
