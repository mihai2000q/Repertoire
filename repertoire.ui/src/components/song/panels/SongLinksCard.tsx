import Song from '../../../types/models/Song.ts'
import { Anchor, Button, Stack, Text } from '@mantine/core'
import { IconBrandYoutube, IconGuitarPick } from '@tabler/icons-react'
import EditPanelCard from '../../@ui/card/EditPanelCard.tsx'
import EditSongLinksModal from '../modal/EditSongLinksModal.tsx'
import { useDisclosure } from '@mantine/hooks'
import YoutubeModal from '../../@ui/modal/YoutubeModal.tsx'
import YoutubeContextMenu from '../../@ui/menu/YoutubeContextMenu.tsx'

interface SongLinksCardProps {
  song: Song
}

function SongLinksCard({ song }: SongLinksCardProps) {
  const [openedEdit, { open: openEdit, close: closeEdit }] = useDisclosure(false)
  const [openedYoutube, { open: openYoutube, close: closeYoutube }] = useDisclosure(false)

  return (
    <EditPanelCard p={'md'} ariaLabel={'song-links-card'} onEditClick={openEdit}>
      <Stack>
        <Text fw={600}>Links</Text>
        <Stack gap={'xs'} align={'center'}>
          {!song.youtubeLink && !song.songsterrLink && (
            <Text fw={500} c={'dimmed'} ta={'center'} fs={'italic'}>
              No links to display
            </Text>
          )}
          {song.songsterrLink && (
            <Anchor
              w={'100%'}
              underline={'never'}
              href={song.songsterrLink}
              target="_blank"
              rel="noreferrer"
              style={{ justifyItems: 'center' }}
            >
              <Button
                fullWidth
                maw={400}
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
          {song.youtubeLink && (
            <YoutubeContextMenu title={song.title} link={song.youtubeLink}>
              <Button
                fullWidth
                maw={400}
                variant={'gradient'}
                size={'md'}
                radius={'lg'}
                leftSection={<IconBrandYoutube size={30} />}
                fz={'h6'}
                fw={700}
                gradient={{ from: 'red.7', to: 'red.1', deg: 90 }}
                onClick={openYoutube}
              >
                Youtube
              </Button>
            </YoutubeContextMenu>
          )}
        </Stack>
      </Stack>

      <EditSongLinksModal song={song} opened={openedEdit} onClose={closeEdit} />
      <YoutubeModal
        title={song.title}
        link={song.youtubeLink}
        opened={openedYoutube}
        onClose={closeYoutube}
      />
    </EditPanelCard>
  )
}

export default SongLinksCard
