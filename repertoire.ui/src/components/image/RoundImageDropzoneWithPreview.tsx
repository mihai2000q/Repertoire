import { Dispatch, SetStateAction } from 'react'
import { ActionIcon, alpha, Box, FileButton, Group, Image, Tooltip } from '@mantine/core'
import { IconPhotoDown, IconUpload, IconUserFilled, IconX } from '@tabler/icons-react'
import { Dropzone, FileWithPath, IMAGE_MIME_TYPE } from '@mantine/dropzone'

interface ImageDropzoneWithPreviewProps {
  image: FileWithPath
  setImage: Dispatch<SetStateAction<FileWithPath>>
}

function ImageDropzoneWithPreview({ image, setImage }: ImageDropzoneWithPreviewProps) {
  if (image) {
    return (
      <Box pos={'relative'}>
        <Image src={URL.createObjectURL(image)} w={92} h={92} radius={'50%'} alt={'song-image'} />

        <Box pos={'absolute'} top={70} left={-8}>
          <Tooltip label={'Remove Image'} openDelay={300} position={'bottom'}>
            <ActionIcon
              c={'white'}
              radius={'50%'}
              aria-label={'remove-image-button'}
              size={'lg'}
              sx={(theme) => ({
                transition: '0.15s',
                color: theme.colors.red[4],
                backgroundColor: alpha(theme.colors.red[5], 0.5),
                '&:hover': {
                  color: theme.colors.red[6],
                  backgroundColor: alpha(theme.colors.red[5], 0.7)
                }
              })}
              onClick={() => setImage(null)}
            >
              <IconX size={18} />
            </ActionIcon>
          </Tooltip>
        </Box>
        <Box pos={'absolute'} top={70} right={-8}>
          <FileButton onChange={setImage} accept="image/png,image/jpeg">
            {(props) => (
              <Tooltip label={'Reload Image'} openDelay={300} position={'bottom'}>
                <ActionIcon
                  c={'dark'}
                  radius={'50%'}
                  aria-label={'add-image-button'}
                  size={'lg'}
                  {...props}
                  sx={(theme) => ({
                    transition: '0.15s',
                    backgroundColor: alpha(theme.colors.gray[4], 0.5),
                    '&:hover': {
                      backgroundColor: alpha(theme.colors.gray[4], 0.7)
                    }
                  })}
                >
                  <IconPhotoDown size={18} />
                </ActionIcon>
              </Tooltip>
            )}
          </FileButton>
        </Box>
      </Box>
    )
  }

  return (
    <Dropzone
      onDrop={(files) => setImage(files[0])}
      accept={IMAGE_MIME_TYPE}
      multiple={false}
      w={92}
      h={92}
      styles={{
        inner: {
          height: 'calc(100% - 5px)'
        }
      }}
      sx={(theme) => ({
        cursor: 'pointer',
        borderRadius: '50%',
        border: `1px solid ${theme.colors.gray[4]}`,
        transition: '0.3s',
        color: theme.colors.gray[6],
        backgroundColor: theme.colors.gray[3],
        '&:hover': {
          color: theme.colors.gray[7],
          backgroundColor: theme.colors.gray[4]
        },

        '&:where([data-accept])': {
          color: theme.colors.green[7],
          backgroundColor: theme.colors.green[3]
        },

        '&:where([data-reject])': {
          color: theme.colors.red[7],
          backgroundColor: theme.colors.red[3]
        }
      })}
    >
      <Group justify="center" style={{ pointerEvents: 'none' }} h={'100%'}>
        <Dropzone.Accept>
          <IconUpload size={45} />
        </Dropzone.Accept>
        <Dropzone.Reject>
          <IconX size={45} />
        </Dropzone.Reject>
        <Dropzone.Idle>
          <IconUserFilled size={45} />
        </Dropzone.Idle>
      </Group>
    </Dropzone>
  )
}

export default ImageDropzoneWithPreview
