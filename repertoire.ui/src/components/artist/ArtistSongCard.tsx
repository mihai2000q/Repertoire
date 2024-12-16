import Song from '../../types/models/Song.ts'
import { ActionIcon, alpha, Avatar, Group, Menu, Stack, Text } from '@mantine/core'
import songPlaceholder from '../../assets/image-placeholder-1.jpg'
import dayjs from 'dayjs'
import { useAppDispatch } from '../../state/store.ts'
import { openAlbumDrawer, openSongDrawer } from '../../state/globalSlice.ts'
import { MouseEvent, useState } from 'react'
import { IconDots, IconTrash } from '@tabler/icons-react'
import { useDisclosure, useHover } from '@mantine/hooks'
import WarningModal from '../@ui/modal/WarningModal.tsx'

interface ArtistSongCardProps {
  song: Song
  handleRemove: () => void
  isUnknownArtist: boolean
}

function ArtistSongCard({ song, handleRemove, isUnknownArtist }: ArtistSongCardProps) {
  const dispatch = useAppDispatch()
  const { ref, hovered } = useHover()

  const [isMenuOpened, setIsMenuOpened] = useState(false)

  const isSelected = hovered || isMenuOpened

  const [openedRemoveWarning, { open: openRemoveWarning, close: closeRemoveWarning }] =
    useDisclosure(false)

  function handleClick() {
    dispatch(openSongDrawer(song.id))
  }

  function handleAlbumClick(e: MouseEvent) {
    e.stopPropagation()
    dispatch(openAlbumDrawer(song.album.id))
  }

  function handleOpenRemoveWarning(e: MouseEvent) {
    e.stopPropagation()
    openRemoveWarning()
  }

  return (
    <>
      <Group
        ref={ref}
        align={'center'}
        wrap={'nowrap'}
        sx={(theme) => ({
          cursor: 'default',
          transition: '0.3s',
          ...(isSelected && {
            boxShadow: theme.shadows.xl,
            backgroundColor: alpha(theme.colors.cyan[0], 0.15)
          })
        })}
        px={'md'}
        py={'xs'}
        onClick={handleClick}
      >
        <Avatar radius={'8px'} src={song.imageUrl ?? songPlaceholder} />

        <Stack gap={0} flex={1} style={{ overflow: 'hidden' }}>
          <Group gap={4}>
            <Text fw={500} truncate={'end'}>
              {song.title}
            </Text>
            {song.album && (
              <>
                <Text fz={'sm'}>-</Text>
                <Text
                  fz={'sm'}
                  c={'dimmed'}
                  truncate={'end'}
                  sx={{ '&:hover': { textDecoration: 'underline' } }}
                  onClick={handleAlbumClick}
                >
                  {song.album.title}
                </Text>
              </>
            )}
          </Group>
          {song.releaseDate && (
            <Text fz={'xs'} c={'dimmed'}>
              {dayjs(song.releaseDate).format('DD MMM YYYY')}
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
            {!isUnknownArtist && (
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
              <Text>from this artist?</Text>
            </Group>
          </Stack>
        }
        onYes={handleRemove}
      />
    </>
  )
}

export default ArtistSongCard
