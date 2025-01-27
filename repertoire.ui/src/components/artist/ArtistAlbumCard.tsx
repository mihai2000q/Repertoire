import Album from '../../types/models/Album.ts'
import { ActionIcon, alpha, Avatar, Group, Menu, Space, Stack, Text } from '@mantine/core'
import albumPlaceholder from '../../assets/image-placeholder-1.jpg'
import dayjs from 'dayjs'
import { useAppDispatch } from '../../state/store.ts'
import { openAlbumDrawer } from '../../state/globalSlice.ts'
import { MouseEvent, useState } from 'react'
import { useDisclosure, useHover } from '@mantine/hooks'
import { IconDots, IconEye, IconTrash } from '@tabler/icons-react'
import WarningModal from '../@ui/modal/WarningModal.tsx'
import { useNavigate } from 'react-router-dom'
import useContextMenu from '../../hooks/useContextMenu.ts'

interface ArtistAlbumCardProps {
  album: Album
  handleRemove: () => void
  isUnknownArtist: boolean
}

function ArtistAlbumCard({ album, handleRemove, isUnknownArtist }: ArtistAlbumCardProps) {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const { ref, hovered } = useHover()

  const [openedMenu, menuDropdownProps, { openMenu, closeMenu }] = useContextMenu()
  const [isMenuOpened, setIsMenuOpened] = useState(false)

  const isSelected = hovered || isMenuOpened

  const [openedRemoveWarning, { open: openRemoveWarning, close: closeRemoveWarning }] =
    useDisclosure(false)

  function handleClick() {
    dispatch(openAlbumDrawer(album.id))
  }

  function handleViewDetails(e: MouseEvent) {
    e.stopPropagation()
    navigate(`/album/${album.id}`)
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
      {!isUnknownArtist && (
        <Menu.Item
          leftSection={<IconTrash size={14} />}
          c={'red.5'}
          onClick={handleOpenRemoveWarning}
        >
          Remove
        </Menu.Item>
      )}
    </>
  )

  return (
    <Menu shadow={'lg'} opened={openedMenu} onClose={closeMenu}>
      <Menu.Target>
        <Group
          ref={ref}
          wrap={'nowrap'}
          aria-label={`album-card-${album.title}`}
          sx={(theme) => ({
            cursor: 'default',
            borderRadius: '8px',
            transition: '0.3s',
            ...(isSelected && {
              boxShadow: theme.shadows.xl,
              backgroundColor: alpha(theme.colors.primary[0], 0.15)
            })
          })}
          px={'md'}
          py={'xs'}
          gap={0}
          onClick={handleClick}
          onContextMenu={openMenu}
        >
          <Avatar radius={'8px'} src={album.imageUrl ?? albumPlaceholder} alt={album.title} />

          <Space ml={'md'} />

          <Stack gap={0} flex={1} style={{ overflow: 'hidden' }}>
            <Text fw={500} truncate={'end'}>
              {album.title}
            </Text>
            {album.releaseDate && (
              <Text fz={'xs'} c={'dimmed'}>
                {dayjs(album.releaseDate).format('D MMM YYYY')}
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
        title={`Remove Album`}
        description={
          <Stack gap={4}>
            <Group gap={4}>
              <Text>Are you sure you want to remove</Text>
              <Text fw={600}>{album.title}</Text>
              <Text>from this artist?</Text>
            </Group>
            <Text fz={'sm'} c={'dimmed'}>
              This action will result in the removal of all related songs to this album.
            </Text>
          </Stack>
        }
        onYes={handleRemove}
      />
    </Menu>
  )
}

export default ArtistAlbumCard
