import { Avatar, AvatarGroup, AvatarGroupProps, AvatarProps, Box, Group, Stack, Text, Tooltip } from '@mantine/core'
import { IconUser } from '@tabler/icons-react'
import { BandMember } from '../../../../../types/models/Artist.ts'

interface BandMembersGroupProps extends AvatarGroupProps {
  bandMembers: BandMember[]
  avatarProps?: AvatarProps & {
    iconSize?: number
  }
}

function BandMembersGroup({ bandMembers, avatarProps, ...props }: BandMembersGroupProps) {
  return (
    <Tooltip
      label={
        <Stack gap={6} p={1}>
          {bandMembers.map((bandMember) => (
            <Group key={bandMember.id} gap={5}>
              <Box bg={'white'} style={{ borderRadius: '50%' }}>
                <Avatar
                  key={bandMember.id}
                  size={'xs'}
                  bd={0}
                  color={bandMember.color}
                  src={bandMember.imageUrl}
                  alt={bandMember.imageUrl && bandMember.name}
                >
                  <IconUser aria-label={`default-icon-${bandMember.name}`} size={12} />
                </Avatar>
              </Box>
              <Text fz={'xs'} c={'gray.2'}>
                {bandMember.name}
              </Text>
            </Group>
          ))}
        </Stack>
      }
      openDelay={300}
    >
      <AvatarGroup {...props}>
        {bandMembers.map((bandMember) => (
          <Avatar
            key={bandMember.id}
            size={'28px'}
            bd={0}
            color={bandMember.color}
            src={bandMember.imageUrl}
            alt={bandMember.imageUrl && bandMember.name}
            {...avatarProps}
          >
            <IconUser aria-label={`default-icon-${bandMember.name}`} size={avatarProps.iconSize ?? 15} />
          </Avatar>
        ))}
      </AvatarGroup>
    </Tooltip>
  )
}

export default BandMembersGroup
