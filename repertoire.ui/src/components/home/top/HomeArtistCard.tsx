import Artist from '../../../types/models/Artist.ts'
import { Avatar, Center, Menu, Stack, Text } from '@mantine/core'
import { useState } from 'react'
import { openArtistDrawer } from '../../../state/slice/globalSlice.ts'
import { useAppDispatch } from '../../../state/store.ts'
import CustomIconUserAlt from '../../@ui/icons/CustomIconUserAlt.tsx'
import { useNavigate } from 'react-router-dom'
import useContextMenu from '../../../hooks/useContextMenu.ts'
import { IconEye } from '@tabler/icons-react'

interface HomeArtistCardProps {
  artist: Artist
}

function HomeArtistCard({ artist }: HomeArtistCardProps) {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()

  const [isImageHovered, setIsImageHovered] = useState(false)
  const [openedMenu, menuDropdownProps, { openMenu, closeMenu }] = useContextMenu()

  const isSelected = isImageHovered || openedMenu

  function handleClick() {
    dispatch(openArtistDrawer(artist.id))
  }

  function handleViewDetails() {
    navigate(`/artist/${artist.id}`)
  }

  return (
    <Stack
      aria-label={`artist-card-${artist.name}`}
      align={'center'}
      gap={'xs'}
      style={{ transition: '0.25s', ...(isSelected && { transform: 'scale(1.05)' }) }}
      w={'max(9vw, 140px)'}
    >
      <Menu shadow={'lg'} opened={openedMenu} onClose={closeMenu}>
        <Menu.Target>
          <Avatar
            onMouseEnter={() => setIsImageHovered(true)}
            onMouseLeave={() => setIsImageHovered(false)}
            size={'max(calc(9vw - 25px), 125px)'}
            src={artist.imageUrl}
            alt={artist.imageUrl && artist.name}
            bg={'gray.0'}
            onClick={handleClick}
            onContextMenu={openMenu}
            sx={(theme) => ({
              cursor: 'pointer',
              transition: '0.25s',
              boxShadow: theme.shadows.xl,
              ...(isSelected && { boxShadow: theme.shadows.xxl })
            })}
          >
            <Center c={'gray.7'}>
              <CustomIconUserAlt aria-label={`default-icon-${artist.name}`} size={58} />
            </Center>
          </Avatar>
        </Menu.Target>

        <Menu.Dropdown {...menuDropdownProps}>
          <Menu.Item leftSection={<IconEye size={14} />} onClick={handleViewDetails}>
            View Details
          </Menu.Item>
        </Menu.Dropdown>
      </Menu>

      <Stack w={'100%'} gap={0} style={{ overflow: 'hidden' }}>
        <Text fw={600} lineClamp={2} ta={'center'}>
          {artist.name}
        </Text>
      </Stack>
    </Stack>
  )
}

export default HomeArtistCard
