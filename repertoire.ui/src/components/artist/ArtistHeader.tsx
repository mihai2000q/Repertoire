import Artist from '../../types/models/Artist.ts'
import { Avatar, Center, Group, Menu, Stack, Text, Title } from '@mantine/core'
import { IconEdit, IconInfoSquareRounded, IconQuestionMark, IconTrash } from '@tabler/icons-react'
import plural from '../../utils/plural.ts'
import HeaderPanelCard from '../@ui/card/HeaderPanelCard.tsx'
import ArtistInfoModal from './modal/ArtistInfoModal.tsx'
import EditArtistHeaderModal from './modal/EditArtistHeaderModal.tsx'
import { useDisclosure } from '@mantine/hooks'
import { useNavigate } from 'react-router-dom'
import ImageModal from '../@ui/modal/ImageModal.tsx'
import lowerTitleFontSize from '../../utils/style/lowerTitleFontSize.ts'
import CustomIconUserAlt from '../@ui/icons/CustomIconUserAlt.tsx'
import AddToPlaylistMenuItem from '../@ui/menu/item/AddToPlaylistMenuItem.tsx'
import DeleteArtistModal from '../@ui/modal/delete/DeleteArtistModal.tsx'
import PerfectRehearsalMenuItem from '../@ui/menu/item/PerfectRehearsalMenuItem.tsx'
import { RefObject } from 'react'

interface ArtistHeaderProps {
  ref: RefObject<HTMLDivElement>
  artist: Artist | undefined
  songsTotalCount: number | undefined
  albumsTotalCount: number | undefined
  isUnknownArtist: boolean
}

function ArtistHeader({
  ref,
  artist,
  songsTotalCount,
  albumsTotalCount,
  isUnknownArtist
}: ArtistHeaderProps) {
  const navigate = useNavigate()

  const [openedImage, { open: openImage, close: closeImage }] = useDisclosure(false)
  const [openedArtistInfo, { open: openArtistInfo, close: closeArtistInfo }] = useDisclosure(false)
  const [openedEdit, { open: openEdit, close: closeEdit }] = useDisclosure(false)
  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  const [openedMenu, { open: openMenu, close: closeMenu }] = useDisclosure(false)

  function onDelete() {
    navigate(`/artists`, { replace: true })
  }

  return (
    <HeaderPanelCard
      ref={ref}
      onEditClick={openEdit}
      menuOpened={openedMenu}
      openMenu={openMenu}
      closeMenu={closeMenu}
      menuDropdown={
        <>
          <Menu.Item leftSection={<IconInfoSquareRounded size={14} />} onClick={openArtistInfo}>
            Info
          </Menu.Item>
          <Menu.Item leftSection={<IconEdit size={14} />} onClick={openEdit}>
            Edit
          </Menu.Item>
          <Menu.Divider />

          <AddToPlaylistMenuItem
            ids={[artist?.id]}
            type={'artists'}
            closeMenu={closeMenu}
            disabled={artist?.songsCount === 0}
          />
          <PerfectRehearsalMenuItem id={artist?.id} closeMenu={closeMenu} type={'artist'} />
          <Menu.Divider />

          <Menu.Item leftSection={<IconTrash size={14} />} c={'red.5'} onClick={openDeleteWarning}>
            Delete
          </Menu.Item>
        </>
      }
      hideIcons={isUnknownArtist}
    >
      <Group wrap={'nowrap'}>
        <Avatar
          src={isUnknownArtist ? null : artist.imageUrl}
          alt={!isUnknownArtist && artist.imageUrl ? artist.name : null}
          size={'max(11vw, 125px)'}
          bg={'white'}
          style={(theme) => ({
            boxShadow: theme.shadows.lg,
            ...(!isUnknownArtist && artist.imageUrl && { cursor: 'pointer' })
          })}
          onClick={!isUnknownArtist && artist.imageUrl ? openImage : undefined}
        >
          <Center c={isUnknownArtist ? 'gray.6' : 'gray.7'}>
            {isUnknownArtist ? (
              <IconQuestionMark
                aria-label={'icon-unknown-artist'}
                strokeWidth={3}
                size={'100%'}
                style={{ padding: '12%' }}
              />
            ) : (
              <CustomIconUserAlt
                aria-label={`default-icon-${artist.name}`}
                size={'100%'}
                style={{ padding: '28%' }}
              />
            )}
          </Center>
        </Avatar>
        <Stack gap={'xxs'}>
          {!isUnknownArtist && (
            <Text fw={500} inline>
              Artist
            </Text>
          )}
          {isUnknownArtist ? (
            <Title order={3} fw={200} fs={'italic'} mb={2} fz={'max(2.5vw, 32px)'}>
              Unknown
            </Title>
          ) : (
            <Title order={1} fw={700} lineClamp={2} fz={lowerTitleFontSize(artist.name)}>
              {artist.name}
            </Title>
          )}
          <Text fw={500} fz={'sm'} c={'dimmed'}>
            {!isUnknownArtist && artist.isBand
              ? artist.bandMembers.length + ` member${plural(artist.bandMembers)} • `
              : ''}
            {albumsTotalCount} album{plural(albumsTotalCount)} • {songsTotalCount} song
            {plural(songsTotalCount)}
          </Text>
        </Stack>
      </Group>

      {!isUnknownArtist && (
        <>
          <ImageModal
            opened={openedImage}
            onClose={closeImage}
            title={artist.name}
            image={artist.imageUrl}
          />

          <ArtistInfoModal opened={openedArtistInfo} onClose={closeArtistInfo} artist={artist} />

          <EditArtistHeaderModal artist={artist} opened={openedEdit} onClose={closeEdit} />

          <DeleteArtistModal
            opened={openedDeleteWarning}
            onClose={closeDeleteWarning}
            artist={artist}
            onDelete={onDelete}
          />
        </>
      )}
    </HeaderPanelCard>
  )
}

export default ArtistHeader
