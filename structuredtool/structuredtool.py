from langchain_core.tools import StructuredTool


def multiply(a: int, b: int) -> int:
    """Multiply two numbers."""
    return a * b


'''
async def amultiply(a: int, b: int) -> int:
    """Multiply two numbers."""
    return a * b
'''


#calculator = StructuredTool.from_function(func=multiply, coroutine=amultiply)
calculator = StructuredTool.from_function(func=multiply)

print(calculator.invoke({"a": 2, "b": 3}))
#print(await calculator.ainvoke({"a": 2, "b": 5}))
