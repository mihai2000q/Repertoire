import Album from '../../../types/models/Album.ts'
import { Avatar, Center, Menu, Stack, Text } from '@mantine/core'
import { useState } from 'react'
import { openAlbumDrawer, openArtistDrawer } from '../../../state/slice/globalSlice.ts'
import { useAppDispatch } from '../../../state/store.ts'
import CustomIconAlbumVinyl from '../../@ui/icons/CustomIconAlbumVinyl.tsx'
import useContextMenu from '../../../hooks/useContextMenu.ts'
import { useNavigate } from 'react-router-dom'
import { IconEye, IconUser } from '@tabler/icons-react'

interface HomeAlbumCardProps {
  album: Album
}

function HomeAlbumCard({ album }: HomeAlbumCardProps) {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()

  const [isImageHovered, setIsImageHovered] = useState(false)
  const [openedMenu, menuDropdownProps, { openMenu, closeMenu }] = useContextMenu()

  const isSelected = isImageHovered || openedMenu

  function handleClick() {
    dispatch(openAlbumDrawer(album.id))
  }

  function handleArtistClick() {
    dispatch(openArtistDrawer(album.artist.id))
  }

  function handleViewDetails() {
    navigate(`/album/${album.id}`)
  }

  function handleViewArtist() {
    navigate(`/artist/${album.artist.id}`)
  }

  return (
    <Stack
      aria-label={`album-card-${album.title}`}
      align={'center'}
      gap={0}
      style={{ transition: '0.25s', ...(isSelected && { transform: 'scale(1.05)' }) }}
      w={'max(10vw, 150px)'}
    >
      <Menu shadow={'lg'} opened={openedMenu} onClose={closeMenu}>
        <Menu.Target>
          <Avatar
            onMouseEnter={() => setIsImageHovered(true)}
            onMouseLeave={() => setIsImageHovered(false)}
            radius={'10%'}
            w={'100%'}
            h={'unset'}
            src={album.imageUrl}
            alt={album.imageUrl && album.title}
            bg={'gray.5'}
            onClick={handleClick}
            onContextMenu={openMenu}
            sx={(theme) => ({
              aspectRatio: 1,
              cursor: 'pointer',
              transition: '0.25s',
              boxShadow: theme.shadows.xl,
              ...(isSelected && { boxShadow: theme.shadows.xxl })
            })}
          >
            <Center c={'white'}>
              <CustomIconAlbumVinyl aria-label={`default-icon-${album.title}`} size={40} />
            </Center>
          </Avatar>
        </Menu.Target>

        <Menu.Dropdown {...menuDropdownProps}>
          <Menu.Item leftSection={<IconEye size={14} />} onClick={handleViewDetails}>
            View Details
          </Menu.Item>
          <Menu.Item
            leftSection={<IconUser size={14} />}
            disabled={!album.artist}
            onClick={handleViewArtist}
          >
            View Artist
          </Menu.Item>
        </Menu.Dropdown>
      </Menu>

      <Stack w={'100%'} pt={'xs'} gap={0} style={{ overflow: 'hidden' }}>
        <Text fw={600} lineClamp={2} ta={'center'}>
          {album.title}
        </Text>
        {album.artist && (
          <Text
            fw={500}
            ta={'center'}
            c={'dimmed'}
            truncate={'end'}
            onClick={handleArtistClick}
            sx={{
              cursor: 'pointer',
              '&:hover': { textDecoration: 'underline' }
            }}
          >
            {album.artist.name}
          </Text>
        )}
      </Stack>
    </Stack>
  )
}

export default HomeAlbumCard
