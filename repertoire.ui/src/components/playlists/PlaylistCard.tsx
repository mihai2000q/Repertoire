import Playlist from '../../types/models/Playlist'
import albumPlaceholder from '../../assets/image-placeholder-1.jpg'
import { AspectRatio, Image, Menu, Stack, Text } from '@mantine/core'
import { useNavigate } from 'react-router-dom'
import { IconTrash } from '@tabler/icons-react'
import { toast } from 'react-toastify'
import { useDeletePlaylistMutation } from '../../state/playlistsApi.ts'
import useContextMenu from '../../hooks/useContextMenu.ts'
import {useState} from "react";

interface PlaylistCardProps {
  playlist: Playlist
}

function PlaylistCard({ playlist }: PlaylistCardProps) {
  const navigate = useNavigate()

  const [deletePlaylistMutation] = useDeletePlaylistMutation()

  const [isImageHovered, setIsImageHovered] = useState(false)
  const [openedMenu, menuDropdownProps, { openMenu, onMenuChange }] = useContextMenu()

  function handleClick() {
    navigate(`/playlist/${playlist.id}`)
  }

  function handleDelete() {
    deletePlaylistMutation(playlist.id)
    toast.success(`${playlist.title} deleted!`)
  }

  return (
    <Stack
      align={'center'}
      gap={0}
      style={{ transition: '0.3s', ...(isImageHovered && { transform: 'scale(1.1)' }) }}
      w={150}
    >
      <Menu shadow={'lg'} opened={openedMenu} onChange={onMenuChange}>
        <Menu.Target>
          <AspectRatio>
            <Image
              onMouseEnter={() => setIsImageHovered(true)}
              onMouseLeave={() => setIsImageHovered(false)}
              radius={'lg'}
              src={playlist.imageUrl}
              fallbackSrc={albumPlaceholder}
              onClick={handleClick}
              onContextMenu={openMenu}
              sx={(theme) => ({
                cursor: 'pointer',
                transition: '0.3s',
                boxShadow: theme.shadows.xxl,
                '&:hover': {
                  boxShadow: theme.shadows.xxl_hover
                }
              })}
            />
          </AspectRatio>
        </Menu.Target>

        <Menu.Dropdown {...menuDropdownProps}>
          <Menu.Item c={'red'} leftSection={<IconTrash size={14} />} onClick={handleDelete}>
            Delete
          </Menu.Item>
        </Menu.Dropdown>
      </Menu>

      <Stack w={'100%'} pt={'xs'} style={{ overflow: 'hidden' }}>
        <Text fw={500} lineClamp={2} ta={'center'}>
          {playlist.title}
        </Text>
      </Stack>
    </Stack>
  )
}

export default PlaylistCard