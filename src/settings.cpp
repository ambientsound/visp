/* vi:set ts=8 sts=8 sw=8:
 *
 * PMS  <<Practical Music Search>>
 * Copyright (C) 2006-2009  Kim Tore Jensen
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 *
 * settings.h - configuration option class
 *
 */


#include "settings.h"
#include "config.h"
#include "pms.h"

using namespace std;

extern Pms *	pms;


/*
 * Constructors and destructors
 */
Setting::Setting()
{
	alias = NULL;
	type = SETTING_TYPE_EINVAL;
	key = "";
}

Options::Options()
{
	colors = NULL;
	reset();
}

Options::~Options()
{
	destroy();
}

void		Options::destroy()
{
	vector<Setting *>::iterator	i;
	vector<Topbarline*>::iterator	j;

	/* Truncate old settings array */
	i = vals.begin();
	while (i++ != vals.end())
		delete *i;
	vals.clear();

	/* Truncate topbar */
	j = topbar.begin();
	while (j++ != topbar.end())
		delete *j;
	topbar.clear();

	if (colors != NULL)
		delete colors;
}


/*
 * Reset to defaults
 */
void		Options::reset()
{
	destroy();

	set("scroll", SETTING_TYPE_SCROLL, "normal");
	set("playmode", SETTING_TYPE_PLAYMODE, "linear");
	set("repeatmode", SETTING_TYPE_REPEATMODE, "none");
	set("columns", SETTING_TYPE_FIELDLIST, "artist track title album length");

	set_long("nextinterval", 5);
	set_long("crossfade", 5);
	set_long("mpd_timeout", 30);
	set_long("repeatonedelay", 1);
	set_long("stopdelay", 1);
	set_long("reconnectdelay", 30);
	set_long("directoryminlen", 30);
	set_long("resetstatus", 3);
	set_long("scrolloff", 0);

	set_bool("debug", false);
	set_bool("addtoreturns", false);
	set_bool("ignorecase", true);
	set_bool("regexsearch", false);
	set_bool("followwindow", false);
	set_bool("followcursor", false);
	set_bool("followplayback", false);
	set_bool("nextafteraction", true);
	set_bool("showtopbar", true);
	set_bool("topbarborders", false);
	set_bool("topbarspace", true);
	set_bool("columnspace", true);
	set_bool("mouse", false);

	set_string("directoryformat", "%artist% - %title%");
	set_string("xtermtitle", "PMS: %ifplaying% %artist% - %title% %else% Not playing %endif%");
	set_string("onplaylistfinish", "");
	set_string("libraryroot", "");
	set_string("startuplist", "playlist");
	set_string("librarysort", "default");
	set_string("albumclass", "artist album date"); //FIXME: implement this

	set_string("status_unknown", "??");
	set_string("status_play", "|>");
	set_string("status_pause", "||");
	set_string("status_stop", "[]");

	/*
	 * Set up option aliases
	 */
	alias("ic", "ignorecase");
	alias("so", "scrolloff");

	//TODO: would be nice to have the commented alteratives default if 
	//Unicode is available
	//status_unknown		= "??"; //?
	//status_play		= "|>"; //▶
	//status_pause		= "||"; //‖
	//status_stop		= "[]"; //■
	
	/* Set up default top bar values */
	topbar.clear();
	while(topbar.size() < 3)
		topbar.push_back(new Topbarline());

	topbar[0]->strings[0] = _("%time_elapsed% %playstate% %time%%ifcursong% (%progresspercentage%%%)%endif%");
	topbar[0]->strings[1] = _("%ifcursong%%artist%%endif%");
	topbar[0]->strings[2] = _("Vol: %volume%%%  Mode: %muteshort%%repeatshort%%randomshort%%manualshort%");
	topbar[1]->strings[1] = _("%ifcursong%==> %title% <==%else%No current song%endif%");
	topbar[2]->strings[0] = _("%listsize%");
	topbar[2]->strings[1] = _("%ifcursong%%album% (%year%)%endif%");
	topbar[2]->strings[2] = _("Q: %livequeuesize%");
}


/*
 * Find a setting based on keyword
 */
Setting *	Options::lookup(string key)
{
	unsigned int	i;

	for (i = 0; i < vals.size(); i++)
	{
		if (vals[i]->key == key)
			return vals[i];
	}

	return NULL;
}


/*
 * Initialize a setting or return an existing one with type t.
 * Returns NULL if the setting exists but has a different SettingType,
 * else return the pointer to the new object.
 */
Setting *	Options::add(string key, SettingType t)
{
	Setting *	s;

	s = lookup(key);
	if (s != NULL)
	{
		if (s->type != t)
			return NULL;
		return s;
	}

	s = new Setting();
	if (s == NULL)
		return s;
	
	s->key = key;
	s->type = t;

	vals.push_back(s);

	return s;
}


/*
 * Alias a keyword to point to another keyword.
 * Limitations: needs to be added _after_ the original keyword.
 * Can be nested indefinately.
 */
bool		Options::alias(string key, string dest)
{
	Setting *	s_key;
	Setting *	s_dest;

	s_dest = lookup(dest);
	if (s_dest == NULL)
		return false;
	
	s_key = add(key, SETTING_TYPE_ALIAS);
	if (s_key == NULL)
		return false;
	
	s_key->alias = s_dest;

	return true;
}


/*
 * Set a variable, make best guess.
 * This function expects options to already exist. Use it from the Configurator class.
 */
bool		Options::set(string key, string val)
{
	Setting *	s;

	if (key.size() > 6 && key.substr(0, 6) == "topbar")
	{
		if (set_topbar_values(key, val))
		{
			set_string(key, val);
			pms->mediator->add("setting.topbar");
			return true;
		}
	}

	s = lookup(key);
	if (s == NULL)
	{
		err.code = CERR_INVALID_OPTION;
		err.str = _("invalid option");
		err.str += " '" + key + "'";
		return false;
	}

	if (set(key, s->type, val) != NULL)
	{
		pms->mediator->add("setting." + key);
		return true;
	}
	return false;
}

/*
 * Set a special or generic variable.
 * This function does all type error checking and converts strings to other types.
 */
Setting *	Options::set(string key, SettingType t, string val)
{
	Setting *	s;

	s = add(key, t);
	if (s == NULL)
		return s;

	switch(t)
	{
		case SETTING_TYPE_STRING:
			return set_string(key, val);

		case SETTING_TYPE_LONG:
			return set_long(key, atoi(val.c_str()));

		case SETTING_TYPE_BOOLEAN:
			return set_bool(key, Configurator::strtobool(val));

		case SETTING_TYPE_FIELDLIST:
			if (Configurator::verify_columns(val, err))
				return set_string(key, val);
			else
				return NULL;

		case SETTING_TYPE_PLAYMODE:
			if (val == "manual")
				set_long(key, PLAYMODE_MANUAL);
			else if (val == "linear")
				set_long(key, PLAYMODE_LINEAR);
			else if (val == "random")
				set_long(key, PLAYMODE_RANDOM);
			else
			{
				err.code = CERR_INVALID_VALUE;
				err.str = _("invalid play mode, expected 'manual', 'linear' or 'random'");
				return NULL;
			}
			s->v_string = val;
			break;

		case SETTING_TYPE_SCROLL:
			if (val == "centered" || val == "centred")
				set_long(key, SCROLL_CENTERED);
			else if (val == "relative")
				set_long(key, SCROLL_RELATIVE);
			else if (val == "normal")
				set_long(key, SCROLL_NORMAL);
			else
			{
				err.code = CERR_INVALID_VALUE;
				err.str = _("invalid scroll mode, expected 'normal', 'centered' or 'relative'");
				return NULL;
			}
			s->v_string = val;
			break;

		default:
			s->v_string = val;
			s->v_long = 0;
			s->v_bool = 0;
	}

	return s;
}

/*
 * Set a string variable
 */
Setting *	Options::set_string(string key, string val)
{
	Setting *	s;

	s = add(key, SETTING_TYPE_STRING);
	if (s != NULL)
		s->v_string = val;

	return s;
}

/*
 * Set a numeric variable
 */
Setting *	Options::set_long(string key, long val)
{
	Setting *	s;

	s = add(key, SETTING_TYPE_LONG);
	if (s != NULL)
		s->v_long = val;

	return s;
}

/*
 * Set a boolean variable
 */
Setting *	Options::set_bool(string key, bool val)
{
	Setting *	s;

	s = add(key, SETTING_TYPE_BOOLEAN);
	if (s != NULL)
		s->v_bool = val;

	return s;
}

/*
 * Toggle a boolean variable
 */
bool		Options::toggle(string key)
{
	Setting *	s;

	s = lookup(key);
	if (s == NULL)
		return false;
	else if (s->type != SETTING_TYPE_BOOLEAN)
		return false;
	
	s->v_bool = !(s->v_bool);
	return true;
}


/*
 * Read functions
 */
string		Options::get_string(string key)
{
	Setting *	s;

	s = lookup(key);
	if (s == NULL)
		return NULL;

	while (s->alias != NULL)
		s = s->alias;
	
	return s->v_string;
}

long		Options::get_long(string key)
{
	Setting *	s;

	s = lookup(key);
	if (s == NULL)
		return NULL;
	
	while (s->alias != NULL)
		s = s->alias;
	
	return s->v_long;
}

bool		Options::get_bool(string key)
{
	Setting *	s;

	s = lookup(key);
	if (s == NULL)
		return NULL;
	
	while (s->alias != NULL)
		s = s->alias;
	
	return s->v_bool;
}


/*
 * Dump an option
 */
string		Options::dump(string key, Error & e)
{
	Setting *	s;

	s = lookup(key);
	if (s == NULL)
	{
		e.code = CERR_INVALID_OPTION;
		e.str = _("invalid option");
		e.str += " '" + key + "'";
		return false;
	}

	return dump(s);
}

string		Options::dump(Setting * s)
{
	string		r = "";

	if (s == NULL)
		return r;

	switch(s->type)
	{
		default:
		case SETTING_TYPE_STRING:
			r = s->key + "=";
			r += s->v_string;
			break;

		case SETTING_TYPE_LONG:
			r = s->key + "=";
			r += Pms::tostring(s->v_long);
			break;

		case SETTING_TYPE_BOOLEAN:
			if (s->v_bool)
				r = "no";
			r += s->key;
			break;
	}

	return r;
}

/*
 * Dump all options to a long string
 */
string		Options::dump_all()
{
	string		output = "";

	unsigned int	i;
	Setting *	s;

	for (i = 0; i < vals.size(); i++)
	{
		s = vals[i];
		output += "set ";
		output += dump(s);
		output += "\n";
	}

	return output;
}

/*
 * Set which fields to show in the topbar.
 */
bool		Options::set_topbar_values(string name, string value)
{
	int		column, row;

	column = atoi(name.substr(6).c_str());

	if (column <= 0)
		return false;
	else if (column < 10)
		name = name.substr(7);
	else if (column < 100)
		name = name.substr(8);
	else
	{
		err.code = CERR_INVALID_TOPBAR_INDEX;
		err.str = _("invalid topbar line");
		err.str += " '" + Pms::tostring(column) + "', ";
		err.str += _("expected range is 1-99");
		return false;
	}

	if (name.size() == 0)
	{
		err.code = CERR_INVALID_TOPBAR_POSITION;
		err.str = _("expected placement after topbar index");
		return false;
	}

	if (name == ".left")
		row = 0;
	else if (name == ".center" || name == ".centre")
		row = 1;
	else if (name == ".right")
		row = 2;
	else
	{
		err.code = CERR_INVALID_TOPBAR_POSITION;
		err.str = _("invalid topbar position");
		err.str += " '" + name.substr(1) + "', ";
		err.str += _("expected one of: left center right");
		return false;
	}

	/* Arguments are now sanitized, should have row = 0-2, and column = 1-99 */

	while (topbar.size() < column)
		topbar.push_back(new Topbarline());

	topbar[column-1]->strings[row] = value;

	return true;
}
