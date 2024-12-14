import Playlist from '../../../types/models/Playlist.ts'
import { Grid, Modal, Stack, Text } from '@mantine/core'
import dayjs from 'dayjs'

interface PlaylistInfoModalProps {
  opened: boolean
  onClose: () => void
  playlist: Playlist
}

function PlaylistInfoModal({ opened, onClose, playlist }: PlaylistInfoModalProps) {
  return (
    <Modal opened={opened} onClose={onClose} title={'Playlist Info'}>
      <Modal.Body px={'xs'} py={0}>
        <Stack>
          <Grid>
            <Grid.Col span={5}>
              <Text fw={500} c={'dimmed'}>
                Created on:
              </Text>
            </Grid.Col>
            <Grid.Col span={7}>
              <Text fw={600}>{dayjs(playlist.createdAt).format('DD MMMM YYYY, HH:mm')}</Text>
            </Grid.Col>

            <Grid.Col span={5}>
              <Text fw={500} c={'dimmed'}>
                Last Modified on:
              </Text>
            </Grid.Col>
            <Grid.Col span={7}>
              <Text fw={600}>{dayjs(playlist.updatedAt).format('DD MMMM YYYY, HH:mm')}</Text>
            </Grid.Col>
          </Grid>
        </Stack>
      </Modal.Body>
    </Modal>
  )
}

export default PlaylistInfoModal
