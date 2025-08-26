import Song from '../../../types/models/Song.ts'
import { Anchor, Button, Stack, Text } from '@mantine/core'
import { IconBrandYoutube, IconGuitarPick } from '@tabler/icons-react'
import EditWidget from '../../@ui/widget/EditWidget.tsx'
import EditSongLinksModal from '../modal/EditSongLinksModal.tsx'
import { useDisclosure } from '@mantine/hooks'
import YoutubeModal from '../../@ui/modal/YoutubeModal.tsx'

interface SongLinksWidgetProps {
  song: Song
}

function SongLinksWidget({ song }: SongLinksWidgetProps) {
  const [openedEdit, { open: openEdit, close: closeEdit }] = useDisclosure(false)
  const [openedYoutube, { open: openYoutube, close: closeYoutube }] = useDisclosure(false)

  return (
    <EditWidget p={'md'} ariaLabel={'links-widget'} onEditClick={openEdit}>
      <Stack gap={'xs'}>
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
    </EditWidget>
  )
}

export default SongLinksWidget
