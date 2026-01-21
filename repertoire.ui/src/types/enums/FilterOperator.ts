enum FilterOperator {
  Equal = '=',
  GreaterThan = '>',
  GreaterThanOrEqual = '>=',
  In = 'IN',
  LessThan = '<',
  LessThanOrEqual = '<=',
  NotEqual = '<>',
  NotEqualVariant = '!=',
  IsNull = 'IS NULL',
  IsNotNull = 'IS NOT NULL',
  PatternMatching = '~*'
}

export default FilterOperator
