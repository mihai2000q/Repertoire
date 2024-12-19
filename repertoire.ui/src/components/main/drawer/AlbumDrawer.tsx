import {
  ActionIcon,
  AspectRatio,
  Avatar,
  Box,
  Divider,
  Grid,
  Group,
  Image,
  Menu,
  Stack,
  Text,
  Title,
  Tooltip
} from '@mantine/core'
import { useDeleteAlbumMutation, useGetAlbumQuery } from '../../../state/albumsApi.ts'
import { useAppSelector } from '../../../state/store.ts'
import AlbumDrawerLoader from '../loader/AlbumDrawerLoader.tsx'
import imagePlaceholder from '../../../assets/image-placeholder-1.jpg'
import songPlaceholder from '../../../assets/image-placeholder-1.jpg'
import RightSideEntityDrawer from '../../@ui/drawer/RightSideEntityDrawer.tsx'
import { IconDotsVertical, IconEye, IconTrash } from '@tabler/icons-react'
import { useDisclosure } from '@mantine/hooks'
import { useState } from 'react'
import { toast } from 'react-toastify'
import { useNavigate } from 'react-router-dom'
import WarningModal from '../../@ui/modal/WarningModal.tsx'
import userPlaceholder from '../../../assets/user-placeholder.jpg'
import dayjs from 'dayjs'
import plural from '../../../utils/plural.ts'

interface AlbumDrawerProps {
  opened: boolean
  onClose: () => void
}

function AlbumDrawer({ opened, onClose }: AlbumDrawerProps) {
  const navigate = useNavigate()

  const albumId = useAppSelector((state) => state.global.albumDrawer.albumId)
  const [deleteAlbumMutation] = useDeleteAlbumMutation()

  const { data: album, isFetching } = useGetAlbumQuery(albumId, { skip: !albumId })

  const [isHovered, setIsHovered] = useState(false)
  const [isMenuOpened, setIsMenuOpened] = useState(false)

  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  function handleViewDetails() {
    onClose()
    navigate(`/album/${album.id}`)
  }

  function handleDelete() {
    deleteAlbumMutation(album.id)
    onClose()
    toast.success(`${album.title} deleted!`)
  }

  if (!album)
    return (
      <RightSideEntityDrawer
        opened={opened}
        onClose={onClose}
        isLoading={true}
        loader={<AlbumDrawerLoader />}
      />
    )

  return (
    <RightSideEntityDrawer
      opened={opened}
      onClose={onClose}
      isLoading={isFetching}
      loader={<AlbumDrawerLoader />}
    >
      <Stack gap={'xs'}>
        <Box
          onMouseEnter={() => setIsHovered(true)}
          onMouseLeave={() => setIsHovered(false)}
          pos={'relative'}
        >
          <AspectRatio ratio={4 / 3}>
            <Image src={album.imageUrl} fallbackSrc={imagePlaceholder} alt={album.title} />
          </AspectRatio>

          <Box pos={'absolute'} top={0} right={0} p={7}>
            <Menu opened={isMenuOpened} onChange={setIsMenuOpened}>
              <Menu.Target>
                <ActionIcon
                  variant={'grey-subtle'}
                  style={{ transition: '0.25s', opacity: isHovered || isMenuOpened ? 1 : 0 }}
                >
                  <IconDotsVertical size={20} />
                </ActionIcon>
              </Menu.Target>

              <Menu.Dropdown>
                <Menu.Item leftSection={<IconEye size={14} />} onClick={handleViewDetails}>
                  View Details
                </Menu.Item>
                <Menu.Item
                  leftSection={<IconTrash size={14} />}
                  c={'red.5'}
                  onClick={openDeleteWarning}
                >
                  Delete
                </Menu.Item>
              </Menu.Dropdown>
            </Menu>
          </Box>
        </Box>

        <Stack px={'md'} pb={'md'} gap={4}>
          <Title order={5} fw={700}>
            {album.title}
          </Title>

          <Group gap={4}>
            {album.artist && (
              <>
                <Group gap={6}>
                  <Avatar size={28} src={album.artist.imageUrl ?? userPlaceholder} />
                  <Text fw={600} fz={'lg'}>
                    {album.artist.name}
                  </Text>
                </Group>
                <Text c={'dimmed'}>•</Text>
              </>
            )}
            {album.releaseDate && (
              <Tooltip
                label={'Released on ' + dayjs(album.releaseDate).format('DD MMMM YYYY')}
                openDelay={200}
                position={'bottom'}
              >
                <Text fw={500}>{dayjs(album.releaseDate).format('YYYY')} •</Text>
              </Tooltip>
            )}
            <Text fw={500} c={'dimmed'}>
              {album.songs.length} song{plural(album.songs)}
            </Text>
          </Group>

          <Divider my={6} />

          <Stack gap={'md'}>
            {album.songs.map((song) => (
              <Grid key={song.id} align={'center'} gutter={'md'} px={'sm'}>
                <Grid.Col span={1}>
                  <Text fw={500} ta={'center'}>
                    {song.albumTrackNo}
                  </Text>
                </Grid.Col>
                <Grid.Col span={1.4}>
                  <Avatar radius={'8px'} size={28} src={song.imageUrl ?? songPlaceholder} />
                </Grid.Col>
                <Grid.Col span={9.6}>
                  <Text fw={500} truncate={'end'}>
                    {song.title}
                  </Text>
                </Grid.Col>
              </Grid>
            ))}
          </Stack>
        </Stack>
      </Stack>

      <WarningModal
        opened={openedDeleteWarning}
        onClose={closeDeleteWarning}
        title={'Delete Album'}
        description={`Are you sure you want to delete this album?`}
        onYes={handleDelete}
      />
    </RightSideEntityDrawer>
  )
}

export default AlbumDrawer
