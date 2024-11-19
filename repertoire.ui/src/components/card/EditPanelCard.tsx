import {ActionIcon, Box, Card, Tooltip} from '@mantine/core'
import { ReactElement } from 'react'
import { IconPencil } from '@tabler/icons-react'
import {useHover} from "@mantine/hooks";

interface EditPanelCardProps {
  children: ReactElement
  onEditClick: () => void,
  p?: string,
}

function EditPanelCard({ children, onEditClick, p }: EditPanelCardProps) {
  const { ref, hovered } = useHover()

  return (
    <Card ref={ref} variant={'panel'} p={p}>
      {children}

      <Box pos={'absolute'} right={0} top={0} p={4}>
        <Tooltip label={'Edit Panel'} openDelay={500}>
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
    </Card>
  )
}

export default EditPanelCard
