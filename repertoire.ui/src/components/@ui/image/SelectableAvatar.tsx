import { forwardRef, MouseEvent } from 'react'
import { alpha, Avatar, AvatarProps, Box } from '@mantine/core'
import { IconCheck } from '@tabler/icons-react'

interface SelectableAvatarProps extends AvatarProps {
  isSelected: boolean
  id?: string
  onClick?: (e: MouseEvent) => void
}

const SelectableAvatar = forwardRef<HTMLDivElement, SelectableAvatarProps>(
  ({ isSelected, id, onClick, ...others }, ref) => {
    return (
      <Box pos={'relative'}>
        <Avatar ref={ref} id={id} onClick={onClick} {...others} />

        <Box
          data-testid={'selected-checkmark'}
          pos={'absolute'}
          top={'37%'}
          left={'36%'}
          p={'lg'}
          style={(theme) => ({
            pointerEvents: 'none',
            borderRadius: '100%',
            backgroundColor: alpha(theme.colors.green[2], 0.95),
            transition: '0.2s',
            zIndex: 2,
            opacity: isSelected ? 1 : 0
          })}
        >
          <IconCheck
            color={'white'}
            style={{
              position: 'absolute',
              top: '22%',
              left: '22%'
            }}
          />
        </Box>
      </Box>
    )
  }
)

SelectableAvatar.displayName = 'SelectableAvatar'

export default SelectableAvatar
