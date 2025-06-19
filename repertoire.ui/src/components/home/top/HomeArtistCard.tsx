import Artist from '../../../types/models/Artist.ts'
import { Avatar, Center, Stack, Text } from '@mantine/core'
import { openArtistDrawer } from '../../../state/slice/globalSlice.ts'
import { useAppDispatch } from '../../../state/store.ts'
import CustomIconUserAlt from '../../@ui/icons/CustomIconUserAlt.tsx'
import { useNavigate } from 'react-router-dom'
import { IconEye } from '@tabler/icons-react'
import { useDisclosure, useHover } from '@mantine/hooks'
import { ContextMenu } from '../../@ui/menu/ContextMenu.tsx'

interface HomeArtistCardProps {
  artist: Artist
}

function HomeArtistCard({ artist }: HomeArtistCardProps) {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const { ref, hovered } = useHover()

  const [openedMenu, { toggle: toggleMenu }] = useDisclosure(false)

  const isSelected = hovered || openedMenu

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
      <ContextMenu shadow={'lg'} opened={openedMenu} onChange={toggleMenu}>
        <ContextMenu.Target>
          <Avatar
            ref={ref}
            size={'max(calc(9vw - 25px), 125px)'}
            src={artist.imageUrl}
            alt={artist.imageUrl && artist.name}
            bg={'gray.0'}
            sx={(theme) => ({
              cursor: 'pointer',
              transition: '0.25s',
              boxShadow: theme.shadows.xl,
              ...(isSelected && { boxShadow: theme.shadows.xxl })
            })}
            onClick={handleClick}
          >
            <Center c={'gray.7'}>
              <CustomIconUserAlt aria-label={`default-icon-${artist.name}`} size={58} />
            </Center>
          </Avatar>
        </ContextMenu.Target>

        <ContextMenu.Dropdown>
          <ContextMenu.Item leftSection={<IconEye size={14} />} onClick={handleViewDetails}>
            View Details
          </ContextMenu.Item>
        </ContextMenu.Dropdown>
      </ContextMenu>

      <Stack w={'100%'} gap={0} style={{ overflow: 'hidden' }}>
        <Text fw={600} lineClamp={2} ta={'center'}>
          {artist.name}
        </Text>
      </Stack>
    </Stack>
  )
}

export default HomeArtistCard
