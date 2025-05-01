import Album from '../../types/models/Album.ts'
import {
  ActionIcon,
  alpha,
  Avatar,
  Center,
  Checkbox,
  Flex,
  Group,
  Menu,
  NumberFormatter,
  Space,
  Stack,
  Text,
  Tooltip
} from '@mantine/core'
import dayjs from 'dayjs'
import { useAppDispatch } from '../../state/store.ts'
import { openAlbumDrawer } from '../../state/slice/globalSlice.ts'
import { MouseEvent, useState } from 'react'
import { useDisclosure, useHover } from '@mantine/hooks'
import { IconCircleMinus, IconDots, IconEye, IconTrash } from '@tabler/icons-react'
import WarningModal from '../@ui/modal/WarningModal.tsx'
import { useNavigate } from 'react-router-dom'
import useContextMenu from '../../hooks/useContextMenu.ts'
import Order from '../../types/Order.ts'
import AlbumProperty from '../../types/enums/AlbumProperty.ts'
import { useRemoveAlbumsFromArtistMutation } from '../../state/api/artistsApi.ts'
import { useDeleteAlbumMutation } from '../../state/api/albumsApi.ts'
import CustomIconAlbumVinyl from '../@ui/icons/CustomIconAlbumVinyl.tsx'
import ConfidenceBar from '../@ui/bar/ConfidenceBar.tsx'
import ProgressBar from '../@ui/bar/ProgressBar.tsx'

interface ArtistAlbumCardProps {
  album: Album
  artistId: string
  isUnknownArtist: boolean
  order: Order
}

function ArtistAlbumCard({ album, artistId, isUnknownArtist, order }: ArtistAlbumCardProps) {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const { ref, hovered } = useHover()

  const [removeAlbumsFromArtist, { isLoading: isRemoveLoading }] =
    useRemoveAlbumsFromArtistMutation()
  const [deleteAlbum, { isLoading: isDeleteLoading }] = useDeleteAlbumMutation()

  const [deleteWithSongs, setDeleteWithSongs] = useState(false)

  const [openedMenu, menuDropdownProps, { openMenu, closeMenu }] = useContextMenu()
  const [isMenuOpened, setIsMenuOpened] = useState(false)

  const isSelected = hovered || isMenuOpened || openedMenu

  const [openedRemoveWarning, { open: openRemoveWarning, close: closeRemoveWarning }] =
    useDisclosure(false)
  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
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

  function handleOpenDeleteWarning(e: MouseEvent) {
    e.stopPropagation()
    openDeleteWarning()
  }

  function handleRemoveFromArtist() {
    removeAlbumsFromArtist({ albumIds: [album.id], id: artistId })
  }

  function handleDelete() {
    deleteAlbum({ id: album.id })
  }

  const menuDropdown = (
    <>
      <Menu.Item leftSection={<IconEye size={14} />} onClick={handleViewDetails}>
        View Details
      </Menu.Item>
      {!isUnknownArtist && (
        <Menu.Item leftSection={<IconCircleMinus size={14} />} onClick={handleOpenRemoveWarning}>
          Remove from artist
        </Menu.Item>
      )}
      <Menu.Item
        leftSection={<IconTrash size={14} />}
        c={'red.5'}
        onClick={handleOpenDeleteWarning}
      >
        Delete
      </Menu.Item>
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
          <Avatar
            radius={'md'}
            src={album.imageUrl}
            alt={album.imageUrl && album.title}
            bg={'gray.5'}
          >
            <Center c={'white'}>
              <CustomIconAlbumVinyl aria-label={`default-icon-${album.title}`} size={15} />
            </Center>
          </Avatar>

          <Space ml={'md'} />

          <Stack gap={0} flex={1} style={{ overflow: 'hidden' }}>
            <Text fw={500} lineClamp={order.property === AlbumProperty.Title ? 2 : 1}>
              {album.title}
            </Text>
            <Flex>
              {order.property === AlbumProperty.ReleaseDate && (
                <Tooltip
                  label={`Album was released on ${dayjs(album.releaseDate).format('D MMMM YYYY')}`}
                  openDelay={400}
                  disabled={!album.releaseDate}
                >
                  <Text fz={'xs'} c={'dimmed'}>
                    {album.releaseDate && dayjs(album.releaseDate).format('D MMM YYYY')}
                  </Text>
                </Tooltip>
              )}
              {order.property === AlbumProperty.Rehearsals && (
                <Tooltip.Floating
                  role={'tooltip'}
                  label={
                    <>
                      Rehearsals: <NumberFormatter value={album.rehearsals} />
                    </>
                  }
                >
                  <Text fz={'xs'} c={'dimmed'}>
                    <NumberFormatter value={album.rehearsals} />
                  </Text>
                </Tooltip.Floating>
              )}
              {order.property === AlbumProperty.Confidence && (
                <ConfidenceBar confidence={album.confidence} w={100} mt={4} />
              )}
              {order.property === AlbumProperty.Progress && (
                <ProgressBar progress={album.progress} w={100} mt={4} />
              )}
              {order.property === AlbumProperty.LastPlayed && (
                <Tooltip
                  label={`Album was played last time on ${dayjs(album.lastTimePlayed).format('D MMMM YYYY [at] hh:mm A')}`}
                  openDelay={400}
                  disabled={!album.lastTimePlayed}
                >
                  <Text fz={'xs'} c={'dimmed'}>
                    {album.lastTimePlayed
                      ? dayjs(album.lastTimePlayed).format('D MMM YYYY')
                      : 'never'}
                  </Text>
                </Tooltip>
              )}
            </Flex>
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
        title={`Remove Album From Artist`}
        description={
          <Stack gap={5}>
            <Group gap={'xxs'}>
              <Text>Are you sure you want to delete</Text>
              <Text fw={600}>{album.title}</Text>
              <Text>?</Text>
            </Group>
            <Checkbox
              checked={deleteWithSongs}
              onChange={(event) => setDeleteWithSongs(event.currentTarget.checked)}
              label={'Delete all associated songs'}
              c={'dimmed'}
              styles={{ label: { paddingLeft: 8 } }}
            />
          </Stack>
        }
        isLoading={isRemoveLoading}
        onYes={handleRemoveFromArtist}
      />
      <WarningModal
        opened={openedDeleteWarning}
        onClose={closeDeleteWarning}
        title={`Delete Album`}
        description={
          <Group gap={'xxs'}>
            <Text>Are you sure you want to delete</Text>
            <Text fw={600}>{album.title}</Text>
            <Text>?</Text>
          </Group>
        }
        isLoading={isDeleteLoading}
        onYes={handleDelete}
      />
    </Menu>
  )
}

export default ArtistAlbumCard
