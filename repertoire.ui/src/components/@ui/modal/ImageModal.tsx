import { alpha, Box, CloseButton, Image, Modal } from '@mantine/core'

interface ImageModalProps {
  opened: boolean
  onClose: () => void
  title: string
  image: string
}

function ImageModal({ opened, onClose, title, image }: ImageModalProps) {
  return (
    <Modal.Root
      role={'dialog'}
      aria-label={title + '-image'}
      opened={opened}
      onClose={onClose}
      size={'auto'}
      centered
    >
      <Modal.Overlay blur={1} backgroundOpacity={0.73} />
      <Modal.Content radius={'lg'}>
        <Modal.Body p={0}>
          <Box pos={'relative'}>
            <Image w={'100%'} h={'max(50vh, 350px)'} src={image} alt={title} />

            <CloseButton
              c={'white'}
              radius={'50%'}
              size={'lg'}
              iconSize={20}
              pos={'absolute'}
              top={0}
              right={0}
              mt={'xs'}
              mr={'md'}
              onClick={onClose}
              sx={(theme) => ({
                transition: '0.25s',
                backgroundColor: alpha(theme.colors.gray[6], 0.25),
                '&:hover': {
                  backgroundColor: alpha(theme.colors.gray[6], 0.5)
                }
              })}
            />
          </Box>
        </Modal.Body>
      </Modal.Content>
    </Modal.Root>
  )
}

export default ImageModal
