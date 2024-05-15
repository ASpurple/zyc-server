export function isNumber(num: unknown) {
	return typeof num === "number" && !isNaN(num);
}

export function isString(value: unknown) {
	return typeof value === "string";
}

export function isBoolean(value: unknown) {
	return typeof value === "boolean";
}

export function isArray(arr: unknown) {
	return arr instanceof Array;
}

export function isStringArray(arr: unknown[]) {
	const isArr = isArray(arr);
	if (!isArr) return false;
	for (let i = 0; i < arr.length; i++) {
		const str = arr[i];
		if (!isString(str)) return false;
	}
	return true;
}

export function isNumberArray(arr: unknown[]) {
	const isArr = isArray(arr);
	if (!isArr) return false;
	for (let i = 0; i < arr.length; i++) {
		const str = arr[i];
		if (!isNumber(str)) return false;
	}
	return true;
}

export function isEmptyValue(val: unknown) {
	return val === undefined || val === null || val === "";
}
