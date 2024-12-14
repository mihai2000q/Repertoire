import Song from '../../../types/models/Song.ts'
import { Anchor, Button, Stack, Text } from '@mantine/core'
import { IconBrandYoutube, IconGuitarPick } from '@tabler/icons-react'
import EditPanelCard from '../../card/EditPanelCard.tsx'
import EditSongLinksModal from '../modal/EditSongLinksModal.tsx'
import { useDisclosure } from '@mantine/hooks'

interface SongLinksCardProps {
  song: Song
}

function SongLinksCard({ song }: SongLinksCardProps) {
  const [openedEdit, { open: openEdit, close: closeEdit }] = useDisclosure(false)

  return (
    <EditPanelCard p={'md'} onEditClick={openEdit}>
      <Stack>
        <Text fw={600}>Links</Text>
        <Stack gap={'xs'}>
          {!song.youtubeLink && !song.songsterrLink && (
            <Text fw={500} c={'dimmed'} ta={'center'} fs={'italic'}>
              No links to display
            </Text>
          )}
          {song.youtubeLink && (
            <Anchor underline={'never'} href={song.youtubeLink} target="_blank" rel="noreferrer">
              <Button
                fullWidth
                variant={'gradient'}
                size={'md'}
                radius={'lg'}
                leftSection={<IconBrandYoutube size={30} />}
                fz={'h6'}
                fw={700}
                gradient={{ from: 'red.7', to: 'red.1', deg: 90 }}
              >
                Youtube
              </Button>
            </Anchor>
          )}
          {song.songsterrLink && (
            <Anchor underline={'never'} href={song.songsterrLink} target="_blank" rel="noreferrer">
              <Button
                fullWidth
                variant={'gradient'}
                size={'md'}
                radius={'lg'}
                leftSection={<IconGuitarPick size={30} />}
                fz={'h6'}
                fw={700}
                gradient={{ from: 'blue.7', to: 'blue.1', deg: 90 }}
              >
                Songsterr
              </Button>
            </Anchor>
          )}
        </Stack>
      </Stack>

      <EditSongLinksModal song={song} opened={openedEdit} onClose={closeEdit} />
    </EditPanelCard>
  )
}

export default SongLinksCard
