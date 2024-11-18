import Artist from '../../types/models/Artist.ts'
import { Avatar, Menu, Stack, Text } from '@mantine/core'
import artistPlaceholder from '../../assets/user-placeholder.jpg'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { IconTrash } from '@tabler/icons-react'
import { toast } from 'react-toastify'
import useContextMenu from '../../hooks/useContextMenu.ts'
import { useDeleteArtistMutation } from '../../state/artistsApi.ts'

interface ArtistCardProps {
  artist: Artist
}

function ArtistCard({ artist }: ArtistCardProps) {
  const navigate = useNavigate()

  const [deleteArtistMutation] = useDeleteArtistMutation()

  const [isAvatarHovered, setIsAvatarHovered] = useState(false)

  const [openedMenu, menuDropdownProps, { openMenu, onMenuChange }] = useContextMenu()

  function handleClick() {
    navigate(`/artist/${artist.id}`)
  }

  function handleDelete() {
    deleteArtistMutation(artist.id)
    toast.success(`${artist.name} deleted!`)
  }

  return (
    <Stack
      align={'center'}
      gap={'xs'}
      style={{
        transition: '0.25s',
        ...(isAvatarHovered && {
          transform: 'scale(1.1)'
        })
      }}
    >
      <Menu shadow={'lg'} opened={openedMenu} onChange={onMenuChange}>
        <Menu.Target>
          <Avatar
            onMouseEnter={() => setIsAvatarHovered(true)}
            onMouseLeave={() => setIsAvatarHovered(false)}
            src={artist.imageUrl ?? artistPlaceholder}
            size={125}
            style={(theme) => ({
              cursor: 'pointer',
              transition: '0.3s',
              boxShadow: isAvatarHovered ? theme.shadows.xxl_hover : theme.shadows.xxl
            })}
            onClick={handleClick}
            onContextMenu={openMenu}
          />
        </Menu.Target>

        <Menu.Dropdown {...menuDropdownProps}>
          <Menu.Item c={'red'} leftSection={<IconTrash size={14} />} onClick={handleDelete}>
            Delete
          </Menu.Item>
        </Menu.Dropdown>
      </Menu>
      <Text fw={600} ta={'center'} lineClamp={2}>
        {artist.name}
      </Text>
    </Stack>
  )
}

export default ArtistCard
