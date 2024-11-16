import { ReactNode } from 'react'
import { ActionIcon, Box, Tooltip } from '@mantine/core'
import { IconPencil } from '@tabler/icons-react'
import { useHover } from '@mantine/hooks'

interface EditHeaderCardProps {
  children: ReactNode
  onEditClick: () => void
}

function EditHeaderCard({ children, onEditClick }: EditHeaderCardProps) {
  const { ref, hovered } = useHover()

  return (
    <Box ref={ref} pos={'relative'}>
      {children}

      <Box pos={'absolute'} right={0} bottom={-12} p={0}>
        <Tooltip label={'Edit Header'} openDelay={500}>
          <ActionIcon
            variant={'grey'}
            size={'md'}
            style={{ transition: '0.25s', opacity: hovered ? 1 : 0 }}
            onClick={onEditClick}
          >
            <IconPencil size={18} />
          </ActionIcon>
        </Tooltip>
      </Box>
    </Box>
  )
}

export default EditHeaderCard
