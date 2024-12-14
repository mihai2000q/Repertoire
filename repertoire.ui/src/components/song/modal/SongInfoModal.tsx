import Song from '../../../types/models/Song.ts'
import { Grid, Modal, Stack, Text } from '@mantine/core'
import dayjs from 'dayjs'

interface SongInfoModalProps {
  opened: boolean
  onClose: () => void
  song: Song
}

function SongInfoModal({ opened, onClose, song }: SongInfoModalProps) {
  return (
    <Modal opened={opened} onClose={onClose} title={'Song Info'}>
      <Modal.Body px={'xs'} py={0}>
        <Stack>
          <Grid>
            <Grid.Col span={5}>
              <Text fw={500} c={'dimmed'}>
                Created on:
              </Text>
            </Grid.Col>
            <Grid.Col span={7}>
              <Text fw={600}>{dayjs(song.createdAt).format('DD MMMM YYYY, HH:mm')}</Text>
            </Grid.Col>

            <Grid.Col span={5}>
              <Text fw={500} c={'dimmed'}>
                Last Modified on:
              </Text>
            </Grid.Col>
            <Grid.Col span={7}>
              <Text fw={600}>{dayjs(song.updatedAt).format('DD MMMM YYYY, HH:mm')}</Text>
            </Grid.Col>
          </Grid>
        </Stack>
      </Modal.Body>
    </Modal>
  )
}

export default SongInfoModal
