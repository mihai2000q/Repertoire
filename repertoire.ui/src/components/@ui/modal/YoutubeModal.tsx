import { Box, Modal } from '@mantine/core'

interface YoutubeModalProps {
  opened: boolean
  onClose: () => void
  title: string
  link?: string | undefined
}

function YoutubeModal({ opened, onClose, title, link }: YoutubeModalProps) {
  return (
    <Modal.Root opened={opened} onClose={onClose} size={'min(80vw, 1000px)'} centered>
      <Modal.Overlay />
      <Modal.Content>
        <Modal.Header>
          <Modal.Title fz={'h6'} c={'dark'} fw={900}>
            {title}
          </Modal.Title>
          <Modal.CloseButton />
        </Modal.Header>
        <Modal.Body>
          <Box
            style={{
              position: 'relative',
              paddingBottom: '56%'
            }}
          >
            <iframe
              width={'100%'}
              height={'100%'}
              src={link?.replace('watch?v=', 'embed/')}
              allow='accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture'
              allowFullScreen
              title='Embedded Youtube'
              style={{
                borderRadius: '16px',
                position: 'absolute',
                border: 'none'
              }}
            />
          </Box>
        </Modal.Body>
      </Modal.Content>
    </Modal.Root>
  )
}

export default YoutubeModal
