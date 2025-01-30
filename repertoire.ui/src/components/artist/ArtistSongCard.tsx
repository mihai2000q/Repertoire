import Song from '../../types/models/Song.ts'
import {
  ActionIcon,
  alpha,
  Avatar,
  Group,
  Menu,
  NumberFormatter,
  Stack,
  Text,
  Tooltip
} from '@mantine/core'
import songPlaceholder from '../../assets/image-placeholder-1.jpg'
import dayjs from 'dayjs'
import { useAppDispatch } from '../../state/store.ts'
import { openAlbumDrawer, openSongDrawer } from '../../state/slice/globalSlice.ts'
import { MouseEvent, useState } from 'react'
import { IconDots, IconEye, IconTrash } from '@tabler/icons-react'
import { useDisclosure, useHover } from '@mantine/hooks'
import WarningModal from '../@ui/modal/WarningModal.tsx'
import Order from '../../types/Order.ts'
import SongProperty from '../../utils/enums/SongProperty.ts'
import DifficultyBar from '../@ui/misc/DifficultyBar.tsx'
import SongConfidenceBar from '../@ui/misc/SongConfidenceBar.tsx'
import SongProgressBar from '../@ui/misc/SongProgressBar.tsx'
import useContextMenu from '../../hooks/useContextMenu.ts'
import { useNavigate } from 'react-router-dom'

interface ArtistSongCardProps {
  song: Song
  handleRemove: () => void
  isUnknownArtist: boolean
  order: Order
}

function ArtistSongCard({ song, handleRemove, isUnknownArtist, order }: ArtistSongCardProps) {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const { ref, hovered } = useHover()

  const [openedMenu, menuDropdownProps, { openMenu, closeMenu }] = useContextMenu()
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
          aria-label={`song-card-${song.title}`}
          wrap={'nowrap'}
          sx={(theme) => ({
            cursor: 'default',
            transition: '0.3s',
            ...(isSelected && {
              boxShadow: theme.shadows.xl,
              backgroundColor: alpha(theme.colors.primary[0], 0.15)
            })
          })}
          px={'md'}
          py={'xs'}
          onClick={handleClick}
          onContextMenu={openMenu}
        >
          <Avatar
            radius={'8px'}
            src={song.imageUrl ?? song.album?.imageUrl ?? songPlaceholder}
            alt={song.title}
          />

          <Stack gap={0} flex={1} style={{ overflow: 'hidden' }}>
            <Group gap={'0px 4px'}>
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
                    inline
                  >
                    {song.album.title}
                  </Text>
                </>
              )}
            </Group>
            {order.property === SongProperty.ReleaseDate && (
              <Tooltip
                label={`Song was released on ${dayjs(song.releaseDate).format('DD MMMM YYYY')}`}
                openDelay={400}
                disabled={!song.releaseDate}
              >
                <Text fz={'xs'} c={'dimmed'}>
                  {song.releaseDate && dayjs(song.releaseDate).format('D MMM YYYY')}
                </Text>
              </Tooltip>
            )}
            {order.property === SongProperty.Difficulty && (
              <DifficultyBar difficulty={song.difficulty} maw={25} mt={4} />
            )}
            {order.property === SongProperty.Rehearsals && (
              <Tooltip.Floating
                role={'tooltip'}
                label={
                  <>
                    Rehearsals: <NumberFormatter value={song.rehearsals} />
                  </>
                }
              >
                <Text fz={'xs'} c={'dimmed'}>
                  <NumberFormatter value={song.rehearsals} />
                </Text>
              </Tooltip.Floating>
            )}
            {order.property === SongProperty.Confidence && (
              <SongConfidenceBar confidence={song.confidence} w={100} mt={4} />
            )}
            {order.property === SongProperty.Progress && (
              <SongProgressBar progress={song.progress} w={100} mt={4} />
            )}
            {order.property === SongProperty.LastTimePlayed && (
              <Tooltip
                label={`Song was played last time on ${dayjs(song.lastTimePlayed).format('DD MMMM YYYY')}`}
                openDelay={400}
                disabled={!song.lastTimePlayed}
              >
                <Text fz={'xs'} c={'dimmed'}>
                  {song.lastTimePlayed ? dayjs(song.lastTimePlayed).format('D MMM YYYY') : 'never'}
                </Text>
              </Tooltip>
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
              <Text>from this artist?</Text>
            </Group>
          </Stack>
        }
        onYes={handleRemove}
      />
    </Menu>
  )
}

export default ArtistSongCard
