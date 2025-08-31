import Album from '../../types/models/Album.ts'
import {
  ActionIcon,
  alpha,
  Avatar,
  Center,
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
import { MouseEvent } from 'react'
import { useDisclosure, useHover } from '@mantine/hooks'
import { IconCircleMinus, IconDots, IconEye, IconTrash } from '@tabler/icons-react'
import WarningModal from '../@ui/modal/WarningModal.tsx'
import { useNavigate } from 'react-router-dom'
import Order from '../../types/Order.ts'
import AlbumProperty from '../../types/enums/properties/AlbumProperty.ts'
import { useRemoveAlbumsFromArtistMutation } from '../../state/api/artistsApi.ts'
import CustomIconAlbumVinyl from '../@ui/icons/CustomIconAlbumVinyl.tsx'
import ConfidenceBar from '../@ui/bar/ConfidenceBar.tsx'
import ProgressBar from '../@ui/bar/ProgressBar.tsx'
import AddToPlaylistMenuItem from '../@ui/menu/item/AddToPlaylistMenuItem.tsx'
import { ContextMenu } from '../@ui/menu/ContextMenu.tsx'
import useDoubleMenu from '../../hooks/useDoubleMenu.ts'
import DeleteAlbumModal from '../@ui/modal/delete/DeleteAlbumModal.tsx'
import { toast } from 'react-toastify'
import PerfectRehearsalMenuItem from '../@ui/menu/item/PerfectRehearsalMenuItem.tsx'

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

  const { openedMenu, toggleMenu, openedContextMenu, toggleContextMenu, closeMenus } =
    useDoubleMenu()

  const isSelected = hovered || openedMenu || openedContextMenu

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

  async function handleRemoveFromArtist() {
    await removeAlbumsFromArtist({ albumIds: [album.id], id: artistId }).unwrap()
    toast.success(`${album.title} removed from artist!`)
  }

  const menuDropdown = (
    <>
      <Menu.Item leftSection={<IconEye size={14} />} onClick={handleViewDetails}>
        View Details
      </Menu.Item>
      <Menu.Divider />

      <AddToPlaylistMenuItem
        ids={[album.id]}
        type={'album'}
        closeMenu={closeMenus}
        disabled={album.songsCount === 0}
      />
      <PerfectRehearsalMenuItem id={album.id} closeMenu={closeMenus} type={'album'} />
      <Menu.Divider />

      {!isUnknownArtist && (
        <Menu.Item leftSection={<IconCircleMinus size={14} />} onClick={handleOpenRemoveWarning}>
          Remove from Artist
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
    <ContextMenu opened={openedContextMenu} onChange={toggleContextMenu}>
      <ContextMenu.Target>
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

          <Space
            ml={{ base: 'xs', xs: 'md', sm: 'xs', betweenSmMd: 'md', md: 'xs', lg: 'md' }}
            style={{ transition: '0.16s' }}
          />

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

          <Menu opened={openedMenu} onChange={toggleMenu}>
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
      </ContextMenu.Target>

      <ContextMenu.Dropdown>{menuDropdown}</ContextMenu.Dropdown>

      <WarningModal
        opened={openedRemoveWarning}
        onClose={closeRemoveWarning}
        title={`Remove Album From Artist`}
        description={
          <Group gap={'xxs'}>
            <Text>Are you sure you want to remove</Text>
            <Text fw={600}>{album.title}</Text>
            <Text>from this artist?</Text>
          </Group>
        }
        isLoading={isRemoveLoading}
        onYes={handleRemoveFromArtist}
      />
      <DeleteAlbumModal
        opened={openedDeleteWarning}
        onClose={closeDeleteWarning}
        album={album}
        withName
      />
    </ContextMenu>
  )
}

export default ArtistAlbumCard
