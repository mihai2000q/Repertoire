import Album from '../../../types/models/Album.ts'
import {Grid, Modal, Stack, Text} from '@mantine/core'
import dayjs from "dayjs";

interface AlbumInfoModalProps {
  opened: boolean
  onClose: () => void
  album: Album
}

function AlbumInfoModal({ opened, onClose, album }: AlbumInfoModalProps) {
  return (
    <Modal opened={opened} onClose={onClose} title={'Album Info'}>
      <Modal.Body px={'xs'} py={0}>
        <Stack>
          <Grid>
            <Grid.Col span={5}>
              <Text fw={500} c={'dimmed'}>
                Created on:
              </Text>
            </Grid.Col>
            <Grid.Col span={7}>
              <Text fw={600}>
                {dayjs(album.createdAt).format('DD MMMM YYYY, HH:mm')}
              </Text>
            </Grid.Col>

            <Grid.Col span={5}>
              <Text fw={500} c={'dimmed'}>
                Last Modified on:
              </Text>
            </Grid.Col>
            <Grid.Col span={7}>
              <Text fw={600}>
                {dayjs(album.updatedAt).format('DD MMMM YYYY, HH:mm')}
              </Text>
            </Grid.Col>
          </Grid>
        </Stack>
      </Modal.Body>
    </Modal>
  )
}

export default AlbumInfoModal
