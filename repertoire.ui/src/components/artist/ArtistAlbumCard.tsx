import Album from '../../types/models/Album.ts'
import { ActionIcon, alpha, Avatar, Group, Menu, Stack, Text } from '@mantine/core'
import albumPlaceholder from '../../assets/image-placeholder-1.jpg'
import dayjs from 'dayjs'
import { useAppDispatch } from '../../state/store.ts'
import { openAlbumDrawer } from '../../state/globalSlice.ts'
import { MouseEvent, useState } from 'react'
import { useHover } from '@mantine/hooks'
import { IconDots, IconTrash } from '@tabler/icons-react'

interface ArtistAlbumCardProps {
  album: Album
  handleRemove: () => void
  isUnknownArtist: boolean
}

function ArtistAlbumCard({ album, handleRemove, isUnknownArtist }: ArtistAlbumCardProps) {
  const dispatch = useAppDispatch()
  const { ref, hovered } = useHover()
  const [isMenuOpened, setIsMenuOpened] = useState(false)
  const isSelected = hovered || isMenuOpened

  function handleClick() {
    dispatch(openAlbumDrawer(album.id))
  }

  function handleDelete(e: MouseEvent) {
    e.stopPropagation()
    handleRemove()
  }

  return (
    <Group
      ref={ref}
      align={'center'}
      wrap={'nowrap'}
      sx={(theme) => ({
        cursor: 'default',
        borderRadius: '8px',
        transition: '0.3s',
        ...isSelected && {
          boxShadow: theme.shadows.xl,
          backgroundColor: alpha(theme.colors.cyan[0], 0.15)
        }
      })}
      px={'md'}
      py={'xs'}
      onClick={handleClick}
    >
      <Avatar radius={'8px'} src={album.imageUrl ?? albumPlaceholder} />

      <Stack gap={0} flex={1} style={{ overflow: 'hidden' }}>
        <Text fw={500} truncate={'end'}>
          {album.title}
        </Text>
        {album.releaseDate && (
          <Text fz={'xs'} c={'dimmed'}>
            {dayjs(album.releaseDate).format('DD MMM YYYY')}
          </Text>
        )}
      </Stack>

      <Menu position={'bottom-end'} opened={isMenuOpened} onChange={setIsMenuOpened}>
        <Menu.Target>
          <ActionIcon
            size={'md'}
            variant={'grey'}
            onClick={(e) => e.stopPropagation()}
            style={{
              transition: '0.3s',
              opacity: isSelected ? 1 : 0
            }}
          >
            <IconDots size={15} />
          </ActionIcon>
        </Menu.Target>
        <Menu.Dropdown>
          {!isUnknownArtist && <Menu.Item leftSection={<IconTrash size={14} />} c={'red.5'} onClick={handleDelete}>
            Remove from Artist
          </Menu.Item>}
        </Menu.Dropdown>
      </Menu>
    </Group>
  )
}

export default ArtistAlbumCard
