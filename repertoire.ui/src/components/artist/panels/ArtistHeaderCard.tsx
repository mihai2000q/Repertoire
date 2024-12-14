import Artist from '../../../types/models/Artist.ts'
import { Avatar, Group, Menu, Stack, Text, Title } from '@mantine/core'
import { IconEdit, IconInfoSquareRounded, IconTrash } from '@tabler/icons-react'
import unknownPlaceholder from '../../../assets/unknown-placeholder.png'
import artistPlaceholder from '../../../assets/user-placeholder.jpg'
import plural from '../../../utils/plural.ts'
import HeaderPanelCard from '../../card/HeaderPanelCard.tsx'
import ArtistInfoModal from '../modal/ArtistInfoModal.tsx'
import EditArtistHeaderModal from '../modal/EditArtistHeaderModal.tsx'
import { useDisclosure } from '@mantine/hooks'
import { toast } from 'react-toastify'
import { useDeleteArtistMutation } from '../../../state/artistsApi.ts'
import { useNavigate } from 'react-router-dom'

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

  const [deleteArtistMutation] = useDeleteArtistMutation()

  const [openedArtistInfo, { open: openArtistInfo, close: closeArtistInfo }] = useDisclosure(false)
  const [openedEditArtistHeader, { open: openEditArtistHeader, close: closeEditArtistHeader }] =
    useDisclosure(false)

  function handleDelete() {
    deleteArtistMutation(artist.id)
    navigate(`/artists`, { replace: true })
    toast.success(`${artist.name} deleted!`)
  }

  return (
    <HeaderPanelCard
      onEditClick={openEditArtistHeader}
      menuDropdown={
        <>
          <Menu.Item leftSection={<IconInfoSquareRounded size={14} />} onClick={openArtistInfo}>
            Info
          </Menu.Item>
          <Menu.Item leftSection={<IconEdit size={14} />} onClick={openEditArtistHeader}>
            Edit
          </Menu.Item>
          <Menu.Item leftSection={<IconTrash size={14} />} c={'red.5'} onClick={handleDelete}>
            Delete
          </Menu.Item>
        </>
      }
      hideIcons={isUnknownArtist}
    >
      <Group>
        <Avatar
          src={isUnknownArtist ? unknownPlaceholder : (artist?.imageUrl ?? artistPlaceholder)}
          size={125}
          style={(theme) => ({
            boxShadow: theme.shadows.md
          })}
        />
        <Stack
          gap={4}
          style={{ ...(!isUnknownArtist && { alignSelf: 'start', paddingTop: '16px' }) }}
        >
          {!isUnknownArtist && (
            <Text fw={500} inline>
              Artist
            </Text>
          )}
          {isUnknownArtist ? (
            <Title order={3} fw={200} fs={'italic'} mb={2}>
              Unknown
            </Title>
          ) : (
            <Title order={1} fw={700}>
              {artist?.name}
            </Title>
          )}
          <Text fw={500} fz={'sm'} c={'dimmed'}>
            {albumsTotalCount} album{plural(albumsTotalCount)} • {songsTotalCount} song
            {plural(songsTotalCount)}
          </Text>
        </Stack>
      </Group>

      {!isUnknownArtist && (
        <>
          <ArtistInfoModal opened={openedArtistInfo} onClose={closeArtistInfo} artist={artist} />

          <EditArtistHeaderModal
            artist={artist}
            opened={openedEditArtistHeader}
            onClose={closeEditArtistHeader}
          />
        </>
      )}
    </HeaderPanelCard>
  )
}

export default ArtistHeaderCard
