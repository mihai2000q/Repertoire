import {
  ActionIcon,
  alpha,
  Avatar,
  Combobox,
  ComboboxProps,
  Group,
  ScrollArea,
  Text,
  Tooltip,
  UnstyledButton,
  useCombobox
} from '@mantine/core'
import { BandMember } from '../../../../types/models/Artist.ts'
import { IconSearch, IconUser } from '@tabler/icons-react'
import { useEffect, useState } from 'react'

interface BandMemberCompactSelectProps extends ComboboxProps {
  bandMember: BandMember | null
  setBandMember: (bandMember: BandMember | null) => void
  bandMembers: BandMember[] | undefined
}

function BandMemberCompactSelect({
  bandMember,
  setBandMember,
  bandMembers,
  ...others
}: BandMemberCompactSelectProps) {
  const combobox = useCombobox({
    onDropdownClose: () => {
      combobox.resetSelectedOption()
      combobox.focusTarget()
      setSearch(value)
    },
    onDropdownOpen: () => {
      combobox.focusSearchInput()
    }
  })

  const [value, setValue] = useState<string>(bandMember?.name ?? '')
  const [search, setSearch] = useState(bandMember?.name ?? '')

  useEffect(() => {
    setValue(bandMember?.name ?? '')
    setSearch(bandMember?.name ?? '')
  }, [bandMember])

  const filteredMembers =
    search.trim() !== ''
      ? bandMembers.filter((member) =>
          member.name.toLowerCase().includes(search.toLowerCase().trim())
        )
      : bandMembers

  const BandMemberOption = ({ member }: { member: BandMember }) => (
    <Combobox.Option
      key={member.id}
      value={member.name}
      aria-label={member.name}
      onClick={() => {
        if (bandMember === member) setBandMember(null)
        else setBandMember(member)
      }}
    >
      <Group gap={'xs'} wrap={'nowrap'}>
        <Avatar
          size={'sm'}
          variant={'light'}
          color={member.color}
          src={member.imageUrl}
          alt={member.name}
        >
          <IconUser size={14} />
        </Avatar>
        <Text inline fw={500} lineClamp={2}>
          {member.name}
        </Text>
      </Group>
    </Combobox.Option>
  )

  function handleSubmit(valueString: string) {
    setValue(value)
    setSearch(valueString)
    combobox.closeDropdown()
  }

  function handleClear() {
    setBandMember(null)
  }

  return (
    <Combobox onOptionSubmit={handleSubmit} store={combobox} withArrow {...others}>
      <Combobox.Target withAriaAttributes={false}>
        {bandMember ? (
          <Tooltip label={`${bandMember.name} is selected`} openDelay={200}>
            <UnstyledButton aria-label={bandMember.name} onClick={() => combobox.toggleDropdown()}>
              <Avatar
                size={24}
                color={bandMember.color}
                src={bandMember.imageUrl}
                sx={{ transition: '0.25s', '&:hover': { filter: 'brightness(0.7)' } }}
              >
                <IconUser size={15} />
              </Avatar>
            </UnstyledButton>
          </Tooltip>
        ) : (
          <Tooltip
            label={
              bandMembers
                ? 'Choose a band member'
                : 'The artist of the song is either not set or not a band'
            }
            ta={'center'}
            openDelay={bandMembers ? 200 : 0}
            maw={190}
            multiline
            position={'top'}
            disabled={combobox.dropdownOpened}
          >
            <ActionIcon
              aria-label={'select-band-member'}
              radius={'50%'}
              sx={(theme) => ({
                transition: '0.25s',
                backgroundColor: 'transparent',
                color: theme.colors.gray[5],
                border: `1px dashed ${theme.colors.gray[5]}`,

                '&:hover': {
                  backgroundColor: theme.colors.gray[1],
                  color: alpha(theme.colors.gray[6], 0.85),
                  borderColor: alpha(theme.colors.gray[6], 0.85)
                },
                '&[data-disabled="true"]': {
                  backgroundColor: 'transparent',
                  color: alpha(theme.colors.gray[4], 0.8),
                  borderColor: alpha(theme.colors.gray[4], 0.8)
                }
              })}
              onClick={() => combobox.toggleDropdown()}
              disabled={bandMembers === undefined}
            >
              <IconUser size={15} />
            </ActionIcon>
          </Tooltip>
        )}
      </Combobox.Target>

      <Combobox.Dropdown miw={180} pt={2} px={'xxs'} pb={0}>
        <Combobox.Search
          size={'xs'}
          px={'xxs'}
          pos={'relative'}
          maxLength={100}
          aria-label={'search'}
          placeholder={'Search'}
          leftSection={<IconSearch size={12} />}
          rightSection={bandMember && <Combobox.ClearButton onClear={handleClear} />}
          value={search}
          onChange={(e) => setSearch(e.currentTarget.value)}
          sx={{
            '.mantine-Input-section': {
              position: 'absolute',
              pointerEvents: 'all'
            },
            '.mantine-Combobox-input': { paddingTop: 1 }
          }}
        />
        <Combobox.Options>
          <ScrollArea.Autosize mah={200} scrollbarSize={5}>
            {bandMembers?.length === 0 ? (
              <Combobox.Empty>Artist has no members</Combobox.Empty>
            ) : filteredMembers?.length === 0 ? (
              <Combobox.Empty>No members found</Combobox.Empty>
            ) : (
              filteredMembers?.map((bandMember) => (
                <BandMemberOption key={bandMember.id} member={bandMember} />
              ))
            )}
          </ScrollArea.Autosize>
        </Combobox.Options>
      </Combobox.Dropdown>
    </Combobox>
  )
}

export default BandMemberCompactSelect
