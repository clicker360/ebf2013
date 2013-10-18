/**
 * GETs or SETs form parameters to/from Object.
 * Based on: http://javascriptmvc.com/docs.html#!jQuery.fn.formParams
 * @author Tom
 *
 * GET:
 *  @param {boolean} convert   if true - strings that represent numbers and booleans will be converted and empty string will not be added to the object.
 *  @return {object}           object representing form fields with values
 *
 * SET:
 *  @param {object} params     object of names with values to apply to the form
 *  @param {boolean} clear     if true - fields which values are undefined (in the passed object) will be cleared
 *  @return {object}           jQuery object representing the form (for chaining)
 *
 * @example
 *  <form>
 *    <input name="foo[bar]" value='2'/>
 *    <input name="foo[ced]" value='4'/>
 *  <form/>
 *
 * $('form').formParams() //-> { foo:{bar:'2', ced: '4'} }
 */

(function ($) {
	'use strict';

	var keyBreaker = /[^\[\]]+/g,
		numberMatcher = /^[\-+]?[0-9]*\.?[0-9]+([eE][\-+]?[0-9]+)?$/,
		isNumber = function (value) {
			if (typeof value === 'number') return true;
			if (typeof value !== 'string') return false;
			return value.match(numberMatcher);
		},
		decodeEntities = function (str) {
			var d = document.createElement('div');
			d.innerHTML = str;
			return d.innerText || d.textContent;
		};


	$.fn.extend({
		formParams: function (params, convert) {
			if (typeof params === 'boolean') { convert = params; params = null; }
			if (params) return this.setParams(params, convert);																	// SET
			else if (this[0].nodeName === 'FORM' && this[0].elements) {															// GET
				return jQuery(jQuery.makeArray(this[0].elements)).getParams(convert);
			}
		},


		setParams: function (params, clear) {																					/*jshint eqeqeq: false*/
			this.find('[name]').each(function () {																				// Find all the inputs
				var name = $(this).attr('name'), value = params[name];

				if (name.indexOf('[') > -1) {																					// if name is object, e.g. user[name], userData[address][street], update value to read this correctly
					var names = name.replace(/\]/g, '').split('['), i = 0, n = null, v = params;
					for (; n = names[i++] ;) if (v[n]) v = v[n]; else { v = undefined; break; }
					value = v;
				}

				if (clear !== true && value === undefined) return;																// if clear==true and no value = clear field, otherwise - leave it as it was
				if (value === null || value === undefined) value = '';															// if no value - clear field

				if (typeof value === 'string' && value.indexOf('&') > -1) value = decodeEntities(value);						// decode html special chars (entities)

				if (this.type === 'radio') this.checked = (this.value == value);
				else if (this.type === 'checkbox') this.checked = value;
				else {
					if ('placeholder' in document.createElement('input')) this.value = value;									// normal browser
					else {																										// manually handle placeholders for specIEl browser
						var el = $(this);
						if (this.value != value && value !== '') el.data('changed', true);
						if (value === '') el.data('changed', false).val(el.attr('placeholder'));
						else this.value = value;
					}
				}
			});
			return this;
		},


		getParams: function (convert) {
			var data = {}, current;
			convert = (convert === undefined ? false : convert);

			this.each(function () {
				var el = this, type = el.type && el.type.toLowerCase();
				if ((type === 'submit') || !el.name || el.disabled)  return;													// if we are submit or disabled - ignore

				var key = el.name, value = $.data(el, 'value') || $.fn.val.call([el]),
					parts = key.match(keyBreaker), lastPart;																	// make an array of values

				if (el.type === 'radio' && !el.checked) return;																	// return only "checked" radio value
				if (el.type === 'checkbox') value = el.checked;																	// convert chekbox to [true | false]

				var $el = $(el);
				if ($el.data('changed') !== true && value === $el.attr('placeholder')) value = '';								// clear placeholder valus for IEs

				if (convert) {
					if (isNumber(value)) {
						var tv = parseFloat(value), cmp = tv + '';
						if (value.indexOf('.') > 0) cmp = tv.toFixed(value.split('.')[1].length);								// convert (string)100.00 to (int)100
						if (cmp === value) value = tv;
					}
					else if (value === 'true') value = true;
					else if (value === 'false') value = false;
					if (value === '') value = undefined;
				}

				current = data;
				for (var i = 0; i < parts.length - 1; i++) {																	// go through and create nested objects
					if (!current[parts[i]]) current[parts[i]] = {};
					current = current[parts[i]];
				}
				lastPart = parts[parts.length - 1];

				if (current[lastPart]) {																						// now we are on the last part, set the value
					if (!$.isArray(current[lastPart])) {
						current[lastPart] = current[lastPart] === undefined ? [] : [current[lastPart]];
					}
					current[lastPart].push(value);
				}
				else if (!current[lastPart]) current[lastPart] = value;
			});
			return data;
		}
	});
})(jQuery);