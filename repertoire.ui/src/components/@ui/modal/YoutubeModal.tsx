import { AspectRatio, Modal, useMatches } from '@mantine/core'

interface YoutubeModalProps {
  opened: boolean
  onClose: () => void
  title: string
  link?: string | undefined
}

function YoutubeModal({ opened, onClose, title, link }: YoutubeModalProps) {
  // Not Recommended usage
  const ratio = useMatches({
    base: 1,
    xs: 4 / 3,
    sm: 16 / 9
  })
  const src = link
    ?.replace('watch?v=', 'embed/')
    .replace('youtube', 'youtube-nocookie')
    .replace(/(www\.)?youtu.be/, 'www.youtube-nocookie.com/embed')

  return (
    <Modal.Root
      opened={opened}
      onClose={onClose}
      size={'min(80vw, 1000px)'}
      trapFocus={false}
      centered
    >
      <Modal.Overlay />
      <Modal.Content>
        <Modal.Header>
          <Modal.Title fz={'h6'} c={'dark'} fw={900}>
            {title}
          </Modal.Title>
          <Modal.CloseButton />
        </Modal.Header>
        <Modal.Body>
          <AspectRatio ratio={ratio}>
            <iframe
              width={'100%'}
              height={'100%'}
              src={src}
              allowFullScreen
              title="Embedded Youtube"
              nonce={'youtube-modal'}
              style={{
                borderRadius: '16px',
                border: 'none'
              }}
            />
          </AspectRatio>
        </Modal.Body>
      </Modal.Content>
    </Modal.Root>
  )
}

export default YoutubeModal
