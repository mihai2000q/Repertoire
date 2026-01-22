import { MouseEvent } from 'react'
import { alpha, Avatar, AvatarProps, Box, Center } from '@mantine/core'
import { IconCheck } from '@tabler/icons-react'

interface SelectableAvatarProps extends AvatarProps {
  isSelected: boolean
  id?: string
  onClick?: (e: MouseEvent) => void
  checkmarkSize?: string | number
}

function SelectableAvatar({
  isSelected,
  id,
  onClick,
  checkmarkSize = '63%',
  ...others
}: SelectableAvatarProps) {
  return (
    <Box pos={'relative'}>
      <Avatar id={id} onClick={onClick} {...others} />

      <Center
        data-testid={'selected-overlay'}
        pos={'absolute'}
        top={0}
        left={0}
        w={'100%'}
        h={'100%'}
        style={(theme) => ({
          pointerEvents: 'none',
          borderRadius: others.radius ?? '100%',
          backgroundColor: alpha(theme.white, 0.3),
          transition: '0.2s',
          zIndex: 2,
          opacity: isSelected ? 1 : 0
        })}
      >
        <Center
          data-testid={'selected-checkmark'}
          w={checkmarkSize}
          h={checkmarkSize}
          style={(theme) => ({
            borderRadius: '100%',
            backgroundColor: alpha(theme.colors.green[2], 0.95)
          })}
        >
          <IconCheck color={'white'} size={'75%'} />
        </Center>
      </Center>
    </Box>
  )
}

export default SelectableAvatar
