import Song from '../../types/models/Song.ts'
import { ActionIcon, alpha, Avatar, Group, Menu, Stack, Text } from '@mantine/core'
import songPlaceholder from '../../assets/image-placeholder-1.jpg'
import { useAppDispatch } from '../../state/store.ts'
import { openSongDrawer } from '../../state/globalSlice.ts'
import { useDisclosure, useHover } from '@mantine/hooks'
import { MouseEvent, useState } from 'react'
import { IconDots, IconTrash } from '@tabler/icons-react'
import WarningModal from '../@ui/modal/WarningModal.tsx'

interface AlbumSongCardProps {
  song: Song
  handleRemove: () => void
  isUnknownAlbum: boolean
}

function AlbumSongCard({ song, handleRemove, isUnknownAlbum }: AlbumSongCardProps) {
  const dispatch = useAppDispatch()
  const { ref, hovered } = useHover()

  const [isMenuOpened, setIsMenuOpened] = useState(false)

  const isSelected = hovered || isMenuOpened

  const [openedRemoveWarning, { open: openRemoveWarning, close: closeRemoveWarning }] =
    useDisclosure(false)

  function handleClick() {
    dispatch(openSongDrawer(song.id))
  }

  function handleOpenRemoveWarning(e: MouseEvent) {
    e.stopPropagation()
    openRemoveWarning()
  }

  return (
    <>
      <Group
        aria-label={`song-card-${song.title}`}
        ref={ref}
        align={'center'}
        wrap={'nowrap'}
        sx={(theme) => ({
          cursor: 'default',
          transition: '0.3s',
          '&:hover': {
            boxShadow: theme.shadows.xl,
            backgroundColor: alpha(theme.colors.cyan[0], 0.15)
          }
        })}
        px={'md'}
        py={'xs'}
        onClick={handleClick}
      >
        {!isUnknownAlbum && (
          <Text fw={500} w={35} ta={'center'}>
            {song.albumTrackNo}
          </Text>
        )}
        <Avatar
          radius={'8px'}
          src={song.imageUrl ?? song.album?.imageUrl ?? songPlaceholder}
          alt={song.title}
        />

        <Stack flex={1} style={{ overflow: 'hidden' }}>
          <Text fw={500} truncate={'end'}>
            {song.title}
          </Text>
        </Stack>

        <Menu position={'bottom-end'} opened={isMenuOpened} onChange={setIsMenuOpened}>
          <Menu.Target>
            <ActionIcon
              aria-label={'more-menu'}
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
            {!isUnknownAlbum && (
              <Menu.Item
                leftSection={<IconTrash size={14} />}
                c={'red.5'}
                onClick={handleOpenRemoveWarning}
              >
                Remove
              </Menu.Item>
            )}
          </Menu.Dropdown>
        </Menu>
      </Group>

      <WarningModal
        opened={openedRemoveWarning}
        onClose={closeRemoveWarning}
        title={`Remove Song`}
        description={
          <Stack gap={4}>
            <Group gap={4}>
              <Text>Are you sure you want to remove</Text>
              <Text fw={600}>{song.title}</Text>
              <Text>from this album?</Text>
            </Group>
          </Stack>
        }
        onYes={handleRemove}
      />
    </>
  )
}

export default AlbumSongCard
