import { Avatar, Box, Modal, Stack, Text, Title } from '@mantine/core'
import { BandMember } from '../../../types/models/Artist.ts'
import { IconUser } from '@tabler/icons-react'

interface AddNewBandMemberModalProps {
  opened: boolean
  onClose: () => void
  bandMember: BandMember
}

function BandMemberDetailsModal({ opened, onClose, bandMember }: AddNewBandMemberModalProps) {
  return (
    <Modal opened={opened} onClose={onClose} size={'xs'} withCloseButton={false} trapFocus={false}>
      <Modal.Body p={0}>
        <Box pos={'relative'}>
          <Modal.CloseButton pos={'absolute'} top={0} right={0} onClick={onClose} />

          <Stack gap={'xs'} align={'center'}>
            <Title ta={'center'} fw={800} order={5} lineClamp={3} px={'lg'}>
              {bandMember.name}
            </Title>

            <Stack align={'center'}>
              <Avatar
                variant={'light'}
                size={175}
                radius={'22px'}
                color={bandMember.color}
                src={bandMember.imageUrl}
                alt={bandMember.name}
                style={(theme) => ({ boxShadow: theme.shadows.lg })}
              >
                <IconUser size={60} />
              </Avatar>

              <Stack align={'center'} gap={'xs'}>
                {bandMember.roles.map((role) => (
                  <Text key={role.id} c={'dimmed'} fw={600} inline>
                    {role.name}
                  </Text>
                ))}
              </Stack>
            </Stack>
          </Stack>
        </Box>
      </Modal.Body>
    </Modal>
  )
}

export default BandMemberDetailsModal
