enum FilterOperator {
  Equal = '=',
  GreaterThan = '>',
  GreaterThanOrEqual = '>=',
  In = 'IN',
  LessThan = '<',
  LessThanOrEqual = '<=',
  NotEqual = '<>',
  IsNull = 'IS NULL',
  IsNotNull = 'IS NOT NULL',
  PatternMatching = '~*'
}

export default FilterOperator;
