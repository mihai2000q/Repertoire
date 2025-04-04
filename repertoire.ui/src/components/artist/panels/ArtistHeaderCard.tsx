import Artist from '../../../types/models/Artist.ts'
import { Avatar, Center, Checkbox, Group, Menu, Stack, Text, Title } from '@mantine/core'
import { IconEdit, IconInfoSquareRounded, IconQuestionMark, IconTrash } from '@tabler/icons-react'
import plural from '../../../utils/plural.ts'
import HeaderPanelCard from '../../@ui/card/HeaderPanelCard.tsx'
import ArtistInfoModal from '../modal/ArtistInfoModal.tsx'
import EditArtistHeaderModal from '../modal/EditArtistHeaderModal.tsx'
import { useDisclosure } from '@mantine/hooks'
import { toast } from 'react-toastify'
import { useDeleteArtistMutation } from '../../../state/api/artistsApi.ts'
import { useNavigate } from 'react-router-dom'
import WarningModal from '../../@ui/modal/WarningModal.tsx'
import ImageModal from '../../@ui/modal/ImageModal.tsx'
import { useState } from 'react'
import lowerTitleFontSize from '../../../utils/lowerTitleFontSize.ts'
import CustomIconUserAlt from '../../@ui/icons/CustomIconUserAlt.tsx'

interface ArtistHeaderCardProps {
  artist: Artist | undefined
  songsTotalCount: number | undefined
  albumsTotalCount: number | undefined
  isUnknownArtist: boolean
}

function ArtistHeaderCard({
  artist,
  songsTotalCount,
  albumsTotalCount,
  isUnknownArtist
}: ArtistHeaderCardProps) {
  const navigate = useNavigate()

  const [deleteArtistMutation, { isLoading: isDeleteLoading }] = useDeleteArtistMutation()

  const [deleteWithAssociations, setDeleteWithAssociations] = useState(false)

  const [openedImage, { open: openImage, close: closeImage }] = useDisclosure(false)
  const [openedArtistInfo, { open: openArtistInfo, close: closeArtistInfo }] = useDisclosure(false)
  const [openedEdit, { open: openEdit, close: closeEdit }] = useDisclosure(false)
  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  async function handleDelete() {
    await deleteArtistMutation({
      id: artist.id,
      withAlbums: deleteWithAssociations,
      withSongs: deleteWithAssociations
    }).unwrap()
    navigate(`/artists`, { replace: true })
    toast.success(`${artist.name} deleted!`)
  }

  return (
    <HeaderPanelCard
      onEditClick={openEdit}
      menuDropdown={
        <>
          <Menu.Item leftSection={<IconInfoSquareRounded size={14} />} onClick={openArtistInfo}>
            Info
          </Menu.Item>
          <Menu.Item leftSection={<IconEdit size={14} />} onClick={openEdit}>
            Edit
          </Menu.Item>
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

          <WarningModal
            opened={openedDeleteWarning}
            onClose={closeDeleteWarning}
            title={'Delete Artist'}
            description={
              <Stack gap={'xs'}>
                <Text fw={500}>Are you sure you want to delete this artist?</Text>
                <Checkbox
                  checked={deleteWithAssociations}
                  onChange={(event) => setDeleteWithAssociations(event.currentTarget.checked)}
                  label={
                    <Text c={'dimmed'}>
                      Delete all associated <b>albums</b> and <b>songs</b>
                    </Text>
                  }
                  styles={{ label: { paddingLeft: 8 } }}
                />
              </Stack>
            }
            onYes={handleDelete}
            isLoading={isDeleteLoading}
          />
        </>
      )}
    </HeaderPanelCard>
  )
}

export default ArtistHeaderCard
