import Artist from '../../../types/models/Artist.ts'
import { Grid, Modal, Stack, Text } from '@mantine/core'
import dayjs from 'dayjs'

interface ArtistInfoModalProps {
  opened: boolean
  onClose: () => void
  artist: Artist
}

function ArtistInfoModal({ opened, onClose, artist }: ArtistInfoModalProps) {
  return (
    <Modal opened={opened} onClose={onClose} title={'Artist Info'}>
      <Modal.Body px={'xs'} py={0}>
        <Stack>
          <Grid>
            <Grid.Col span={5}>
              <Text fw={500} c={'dimmed'}>
                Created on:
              </Text>
            </Grid.Col>
            <Grid.Col span={7}>
              <Text fw={600}>{dayjs(artist.createdAt).format('DD MMMM YYYY, HH:mm')}</Text>
            </Grid.Col>

            <Grid.Col span={5}>
              <Text fw={500} c={'dimmed'}>
                Last Modified on:
              </Text>
            </Grid.Col>
            <Grid.Col span={7}>
              <Text fw={600}>{dayjs(artist.updatedAt).format('DD MMMM YYYY, HH:mm')}</Text>
            </Grid.Col>
          </Grid>
        </Stack>
      </Modal.Body>
    </Modal>
  )
}

export default ArtistInfoModal
