import Song from '../../types/models/Song.ts'
import { ActionIcon, alpha, Avatar, Group, Menu, Stack, Text } from '@mantine/core'
import songPlaceholder from '../../assets/image-placeholder-1.jpg'
import { useAppDispatch } from '../../state/store.ts'
import { openAlbumDrawer, openArtistDrawer, openSongDrawer } from '../../state/globalSlice.ts'
import { useDisclosure, useHover } from '@mantine/hooks'
import { MouseEvent, useState } from 'react'
import { IconDots, IconEye, IconTrash } from '@tabler/icons-react'
import WarningModal from '../@ui/modal/WarningModal.tsx'
import { useNavigate } from 'react-router-dom'
import useContextMenu from '../../hooks/useContextMenu.ts'

interface PlaylistSongCardProps {
  song: Song
  handleRemove: () => void
}

function PlaylistSongCard({ song, handleRemove }: PlaylistSongCardProps) {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const { ref, hovered } = useHover()

  const [openedMenu, menuDropdownProps, { openMenu, onMenuChange }] = useContextMenu()
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

  function handleArtistClick(e: MouseEvent) {
    e.stopPropagation()
    dispatch(openArtistDrawer(song.artist.id))
  }

  function handleViewDetails(e: MouseEvent) {
    e.stopPropagation()
    navigate(`/song/${song.id}`)
  }

  function handleOpenRemoveWarning(e: MouseEvent) {
    e.stopPropagation()
    openRemoveWarning()
  }

  const menuDropdown = (
    <>
      <Menu.Item leftSection={<IconEye size={14} />} onClick={handleViewDetails}>
        View Details
      </Menu.Item>
      <Menu.Item
        leftSection={<IconTrash size={14} />}
        c={'red.5'}
        onClick={handleOpenRemoveWarning}
      >
        Remove
      </Menu.Item>
    </>
  )

  return (
    <Menu shadow={'lg'} opened={openedMenu} onChange={onMenuChange}>
      <Menu.Target>
        <Group
          ref={ref}
          wrap={'nowrap'}
          aria-label={`song-card-${song.title}`}
          sx={(theme) => ({
            cursor: 'default',
            transition: '0.3s',
            '&:hover': {
              boxShadow: theme.shadows.xl,
              backgroundColor: alpha(theme.colors.primary[0], 0.15)
            }
          })}
          px={'md'}
          py={'xs'}
          onClick={handleClick}
          onContextMenu={openMenu}
        >
          <Text fw={500} w={35} ta={'center'}>
            {song.playlistTrackNo}
          </Text>

          <Avatar
            radius={'8px'}
            src={song.imageUrl ?? song.album?.imageUrl ?? songPlaceholder}
            alt={song.title}
          />

          <Stack flex={1} gap={0} style={{ overflow: 'hidden' }}>
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
                    style={{ cursor: 'pointer' }}
                    onClick={handleAlbumClick}
                  >
                    {song.album.title}
                  </Text>
                </>
              )}
            </Group>
            {song.artist && (
              <Text
                fz={'sm'}
                c={'dimmed'}
                sx={{ '&:hover': { textDecoration: 'underline' } }}
                style={{ cursor: 'pointer', alignSelf: 'start' }}
                onClick={handleArtistClick}
              >
                {song.artist.name}
              </Text>
            )}
          </Stack>

          <Menu position={'bottom-end'} opened={isMenuOpened} onChange={setIsMenuOpened}>
            <Menu.Target>
              <ActionIcon
                size={'md'}
                variant={'grey'}
                aria-label={'more-menu'}
                onClick={(e) => e.stopPropagation()}
                style={{
                  transition: '0.3s',
                  opacity: isSelected ? 1 : 0
                }}
              >
                <IconDots size={15} />
              </ActionIcon>
            </Menu.Target>

            <Menu.Dropdown>{menuDropdown}</Menu.Dropdown>
          </Menu>
        </Group>
      </Menu.Target>

      <Menu.Dropdown {...menuDropdownProps}>{menuDropdown}</Menu.Dropdown>

      <WarningModal
        opened={openedRemoveWarning}
        onClose={closeRemoveWarning}
        title={`Remove Song`}
        description={
          <Stack gap={4}>
            <Group gap={4}>
              <Text>Are you sure you want to remove</Text>
              <Text fw={600}>{song.title}</Text>
              <Text>from this playlist?</Text>
            </Group>
          </Stack>
        }
        onYes={handleRemove}
      />
    </Menu>
  )
}

export default PlaylistSongCard
