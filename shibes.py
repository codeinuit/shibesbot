import discord
from discord.ext.commands import Bot
from discord.ext import commands

import json
import requests
import os
os.spawnl(os.P_NOWAIT, 'shibesbot')

""" 
"" Discord configuration
"""

Client = discord.Client()
bot_prefix = "shibes"
client = commands.Bot(command_prefix=bot_prefix)
token = "NDM4NzcxMzI2Mzg4MjczMTUz.DcJedQ.HnyD8BMdnmIjhaGgWwJEStuHDvQ-Q"

async def on_shibes_ask(message):
    response = requests.get('http://shibe.online/api/shibes')
    await client.send_message(message.channel, response.json()[0])


@client.event
async def on_message(message):
    if message.content.startswith(bot_prefix):
        await on_shibes_ask(message)


client.run("NDM4NzcxMzI2Mzg4MjczMTUz.DcJfAQ.K9sAt-bd9pQWu49i9MC4Eo7lo1k")