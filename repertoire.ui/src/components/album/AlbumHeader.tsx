import Album from '../../types/models/Album.ts'
import { Avatar, Center, Group, Menu, Stack, Text, Title, Tooltip } from '@mantine/core'
import { IconEdit, IconInfoSquareRounded, IconQuestionMark, IconTrash } from '@tabler/icons-react'
import dayjs from 'dayjs'
import plural from '../../utils/plural.ts'
import HeaderPanelCard from '../@ui/card/HeaderPanelCard.tsx'
import { openArtistDrawer } from '../../state/slice/globalSlice.ts'
import { useDisclosure } from '@mantine/hooks'
import { useAppDispatch } from '../../state/store.ts'
import { useNavigate } from 'react-router-dom'
import AlbumInfoModal from './modal/AlbumInfoModal.tsx'
import EditAlbumHeaderModal from './modal/EditAlbumHeaderModal.tsx'
import ImageModal from '../@ui/modal/ImageModal.tsx'
import titleFontSize from '../../utils/style/titleFontSize.ts'
import CustomIconAlbumVinyl from '../@ui/icons/CustomIconAlbumVinyl.tsx'
import CustomIconUserAlt from '../@ui/icons/CustomIconUserAlt.tsx'
import AddToPlaylistMenuItem from '../@ui/menu/item/AddToPlaylistMenuItem.tsx'
import DeleteAlbumModal from '../@ui/modal/delete/DeleteAlbumModal.tsx'
import PerfectRehearsalMenuItem from '../@ui/menu/item/PerfectRehearsalMenuItem.tsx'

interface AlbumHeaderProps {
  album: Album | undefined
  isUnknownAlbum: boolean
  songsTotalCount: number | undefined
}

function AlbumHeader({ album, isUnknownAlbum, songsTotalCount }: AlbumHeaderProps) {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()

  const [openedImage, { open: openImage, close: closeImage }] = useDisclosure(false)
  const [openedAlbumInfo, { open: openAlbumInfo, close: closeAlbumInfo }] = useDisclosure(false)
  const [openedEdit, { open: openEdit, close: closeEdit }] = useDisclosure(false)
  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  const [openedMenu, { open: openMenu, close: closeMenu }] = useDisclosure(false)

  function handleArtistClick() {
    dispatch(openArtistDrawer(album.artist.id))
  }

  function onDelete() {
    navigate(`/albums`, { replace: true })
  }

  return (
    <HeaderPanelCard
      onEditClick={openEdit}
      menuOpened={openedMenu}
      openMenu={openMenu}
      closeMenu={closeMenu}
      menuDropdown={
        <>
          <Menu.Item leftSection={<IconInfoSquareRounded size={14} />} onClick={openAlbumInfo}>
            Info
          </Menu.Item>
          <Menu.Item leftSection={<IconEdit size={14} />} onClick={openEdit}>
            Edit
          </Menu.Item>
          <Menu.Divider />

          <AddToPlaylistMenuItem
            ids={[album?.id]}
            type={'albums'}
            closeMenu={closeMenu}
            disabled={album?.songsCount === 0}
          />
          <PerfectRehearsalMenuItem id={album?.id} closeMenu={closeMenu} type={'album'} />
          <Menu.Divider />

          <Menu.Item leftSection={<IconTrash size={14} />} c={'red.5'} onClick={openDeleteWarning}>
            Delete
          </Menu.Item>
        </>
      }
      hideIcons={isUnknownAlbum}
    >
      <Group wrap={'nowrap'}>
        <Avatar
          src={isUnknownAlbum ? null : album.imageUrl}
          alt={!isUnknownAlbum && album.imageUrl ? album.title : null}
          radius={'10%'}
          size={'max(12vw, 150px)'}
          bg={isUnknownAlbum ? 'white' : 'gray.5'}
          style={(theme) => ({
            aspectRatio: 1,
            boxShadow: theme.shadows.lg,
            ...(!isUnknownAlbum && album.imageUrl && { cursor: 'pointer' })
          })}
          onClick={!isUnknownAlbum && album.imageUrl ? openImage : undefined}
        >
          <Center c={isUnknownAlbum ? 'gray.6' : 'white'}>
            {isUnknownAlbum ? (
              <IconQuestionMark
                aria-label={'icon-unknown-album'}
                strokeWidth={3}
                size={'100%'}
                style={{ padding: '12%' }}
              />
            ) : (
              <CustomIconAlbumVinyl
                aria-label={`default-icon-${album.title}`}
                size={'100%'}
                style={{ padding: '33%' }}
              />
            )}
          </Center>
        </Avatar>
        <Stack gap={'xxs'}>
          {!isUnknownAlbum && (
            <Text fw={500} inline>
              Album
            </Text>
          )}
          {isUnknownAlbum ? (
            <Title order={3} fw={200} fs={'italic'} fz={'max(2.5vw, 32px)'}>
              Unknown
            </Title>
          ) : (
            <Title order={1} fw={700} lineClamp={2} fz={titleFontSize(album.title)}>
              {album.title}
            </Title>
          )}
          <Group gap={'xxs'} wrap={'nowrap'}>
            {album?.artist && (
              <>
                <Group gap={'xs'} wrap={'nowrap'}>
                  <Avatar
                    size={35}
                    src={album.artist.imageUrl}
                    alt={album.artist.imageUrl && album.artist.name}
                    style={(theme) => ({ boxShadow: theme.shadows.sm })}
                    bg={'gray.0'}
                  >
                    <Center c={'gray.7'}>
                      <CustomIconUserAlt
                        size={15}
                        aria-label={`default-icon-${album.artist.name}`}
                      />
                    </Center>
                  </Avatar>
                  <Text
                    fw={600}
                    fz={'lg'}
                    sx={{
                      cursor: 'pointer',
                      '&:hover': { textDecoration: 'underline' }
                    }}
                    onClick={handleArtistClick}
                    lineClamp={1}
                  >
                    {album.artist.name}
                  </Text>
                </Group>
                <Text c={'dimmed'}>•</Text>
              </>
            )}
            {album?.releaseDate && (
              <>
                <Tooltip
                  label={'Released on ' + dayjs(album.releaseDate).format('D MMMM YYYY')}
                  openDelay={200}
                  position={'bottom'}
                >
                  <Text fw={500} c={'dimmed'}>
                    {dayjs(album.releaseDate).format('YYYY')}
                  </Text>
                </Tooltip>
                <Text fw={500} c={'dimmed'}>
                  •
                </Text>
              </>
            )}
            <Text fw={500} c={'dimmed'} truncate={'end'}>
              {isUnknownAlbum ? songsTotalCount : album.songs.length} song
              {plural(isUnknownAlbum ? songsTotalCount : album.songs)}
            </Text>
          </Group>
        </Stack>
      </Group>

      {!isUnknownAlbum && (
        <>
          <ImageModal
            opened={openedImage}
            onClose={closeImage}
            title={album.title}
            image={album.imageUrl}
          />

          <AlbumInfoModal opened={openedAlbumInfo} onClose={closeAlbumInfo} album={album} />

          <EditAlbumHeaderModal album={album} opened={openedEdit} onClose={closeEdit} />

          <DeleteAlbumModal
            opened={openedDeleteWarning}
            onClose={closeDeleteWarning}
            album={album}
            onDelete={onDelete}
          />
        </>
      )}
    </HeaderPanelCard>
  )
}

export default AlbumHeader
